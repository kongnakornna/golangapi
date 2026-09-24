package vectordb

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"icmongolang/config"
	"icmongolang/pkg/logger"

	"gorm.io/gorm"
)

// pgRow mirrors the vector_documents table row.
type pgRow struct {
	ID         int64     `gorm:"column:id"`
	DocumentID string    `gorm:"column:document_id"`
	Content    string    `gorm:"column:content"`
	SourceType string    `gorm:"column:source_type"`
	Metadata   string    `gorm:"column:metadata"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

type pgScoreRow struct {
	DocumentID string    `gorm:"column:document_id"`
	Content    string    `gorm:"column:content"`
	SourceType string    `gorm:"column:source_type"`
	Metadata   string    `gorm:"column:metadata"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	Score      float64   `gorm:"column:score"`
}

// pgBackend implements VectorDB on top of PostgreSQL + pgvector.
type pgBackend struct {
	table string
	dims  int
	db    *gorm.DB
	log   logger.Logger
}

func newPgBackend(cfg *config.VectorDBConfig, db *gorm.DB, log logger.Logger) VectorDB {
	table := cfg.Index
	if table == "" {
		table = "vector_documents"
	}
	return &pgBackend{table: table, dims: cfg.Dims, db: db, log: log}
}

func (b *pgBackend) Provider() string { return "pgvector" }

func (b *pgBackend) Ping(ctx context.Context) error {
	sqlDB, err := b.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// EnsureIndex enables the pgvector extension and creates the table + HNSW index.
func (b *pgBackend) EnsureIndex(ctx context.Context, dims int) error {
	if dims <= 0 {
		dims = b.dims
	}
	if dims <= 0 {
		dims = 768
	}

	if err := b.db.WithContext(ctx).Exec("CREATE EXTENSION IF NOT EXISTS vector").Error; err != nil {
		return fmt.Errorf("pgvector: enable extension: %w", err)
	}

	createTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id          BIGSERIAL PRIMARY KEY,
		document_id VARCHAR(255) NOT NULL,
		content     TEXT NOT NULL,
		embedding   vector(%d),
		source_type VARCHAR(100) NOT NULL DEFAULT '',
		metadata    TEXT NOT NULL DEFAULT '',
		created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`, b.table, dims)
	if err := b.db.WithContext(ctx).Exec(createTable).Error; err != nil {
		return fmt.Errorf("pgvector: create table: %w", err)
	}

	idxName := b.table + "_embedding_hnsw"
	createIdx := fmt.Sprintf(
		`CREATE INDEX IF NOT EXISTS %s ON %s USING hnsw (embedding vector_cosine_ops)`,
		idxName, b.table,
	)
	if err := b.db.WithContext(ctx).Exec(createIdx).Error; err != nil {
		return fmt.Errorf("pgvector: create hnsw index: %w", err)
	}

	b.log.Infof("✅ pgvector table %q ready (dims=%d)", b.table, dims)
	return nil
}

func (b *pgBackend) Index(ctx context.Context, doc *VectorDoc) (string, error) {
	if doc.Timestamp.IsZero() {
		doc.Timestamp = time.Now()
	}
	vec := vecToString(doc.Embedding)

	var id int64
	insertSQL := fmt.Sprintf(`INSERT INTO %s
		(document_id, content, embedding, source_type, metadata, created_at, updated_at)
		VALUES (?, ?, ?::vector, ?, ?, ?, ?) RETURNING id`, b.table)
	if err := b.db.WithContext(ctx).
		Raw(insertSQL, doc.DocumentID, doc.Content, vec, doc.SourceType, doc.Metadata, doc.Timestamp, doc.Timestamp).
		Scan(&id).Error; err != nil {
		return "", err
	}

	doc.ID = fmt.Sprintf("%d", id)
	return doc.ID, nil
}

func (b *pgBackend) Get(ctx context.Context, id string) (*VectorDoc, error) {
	if !isNumericID(id) {
		return nil, ErrNotFound
	}
	var row pgRow
	selectSQL := fmt.Sprintf(`SELECT id, document_id, content, source_type, metadata, created_at
		FROM %s WHERE id = ?`, b.table)
	if err := b.db.WithContext(ctx).Raw(selectSQL, id).Scan(&row).Error; err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, ErrNotFound
	}
	return &VectorDoc{
		ID:         fmt.Sprintf("%d", row.ID),
		DocumentID: row.DocumentID,
		Content:    row.Content,
		SourceType: row.SourceType,
		Metadata:   row.Metadata,
		Timestamp:  row.CreatedAt,
	}, nil
}

func (b *pgBackend) List(ctx context.Context, limit, offset int) ([]VectorDoc, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	var total int64
	countSQL := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, b.table)
	if err := b.db.WithContext(ctx).Raw(countSQL).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []pgRow
	listSQL := fmt.Sprintf(`SELECT id, document_id, content, source_type, metadata, created_at
		FROM %s ORDER BY id DESC LIMIT ? OFFSET ?`, b.table)
	if err := b.db.WithContext(ctx).Raw(listSQL, limit, offset).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	docs := make([]VectorDoc, 0, len(rows))
	for _, r := range rows {
		docs = append(docs, VectorDoc{
			ID:         fmt.Sprintf("%d", r.ID),
			DocumentID: r.DocumentID,
			Content:    r.Content,
			SourceType: r.SourceType,
			Metadata:   r.Metadata,
			Timestamp:  r.CreatedAt,
		})
	}
	return docs, total, nil
}

// Update merges the non-empty fields of an existing row (partial update).
// Empty fields — including an empty embedding — are left untouched.
func (b *pgBackend) Update(ctx context.Context, id string, doc *VectorDoc) error {
	if !isNumericID(id) {
		return ErrNotFound
	}
	if doc == nil {
		return fmt.Errorf("pgvector: update: doc is nil")
	}

	sets := []string{"updated_at = ?"}
	args := []interface{}{time.Now()}
	if doc.DocumentID != "" {
		sets = append(sets, "document_id = ?")
		args = append(args, doc.DocumentID)
	}
	if doc.Content != "" {
		sets = append(sets, "content = ?")
		args = append(args, doc.Content)
	}
	if len(doc.Embedding) > 0 {
		sets = append(sets, "embedding = ?::vector")
		args = append(args, vecToString(doc.Embedding))
	}
	if doc.SourceType != "" {
		sets = append(sets, "source_type = ?")
		args = append(args, doc.SourceType)
	}
	if doc.Metadata != "" {
		sets = append(sets, "metadata = ?")
		args = append(args, doc.Metadata)
	}
	args = append(args, id)

	updateSQL := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", b.table, strings.Join(sets, ", "))
	res := b.db.WithContext(ctx).Exec(updateSQL, args...)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (b *pgBackend) Delete(ctx context.Context, id string) error {
	if !isNumericID(id) {
		return ErrNotFound
	}
	deleteSQL := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, b.table)
	res := b.db.WithContext(ctx).Exec(deleteSQL, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (b *pgBackend) SearchKNN(ctx context.Context, embedding []float32, k int) ([]SearchHit, error) {
	if k <= 0 {
		k = 10
	}
	vec := vecToString(embedding)

	var rows []pgScoreRow
	searchSQL := fmt.Sprintf(`SELECT document_id, content, source_type, metadata, created_at,
		1 - (embedding <=> ?::vector) AS score
		FROM %s
		WHERE embedding IS NOT NULL
		ORDER BY embedding <=> ?::vector
		LIMIT ?`, b.table)
	if err := b.db.WithContext(ctx).Raw(searchSQL, vec, vec, k).Scan(&rows).Error; err != nil {
		return nil, err
	}

	hits := make([]SearchHit, 0, len(rows))
	for _, r := range rows {
		hits = append(hits, SearchHit{
			DocumentID: r.DocumentID,
			Content:    r.Content,
			SourceType: r.SourceType,
			Metadata:   r.Metadata,
			Timestamp:  r.CreatedAt,
			Score:      r.Score,
		})
	}
	return hits, nil
}

// isNumericID reports whether id is a positive integer (pgvector table PK).
func isNumericID(id string) bool {
	if id == "" {
		return false
	}
	n, err := strconv.ParseInt(id, 10, 64)
	return err == nil && n > 0
}

// vecToString serializes a float32 slice to the PostgreSQL vector literal
// format, e.g. [0.1,0.2,0.3].
func vecToString(v []float32) string {
	if len(v) == 0 {
		return "[]"
	}
	var b strings.Builder
	b.WriteByte('[')
	for i, f := range v {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.FormatFloat(float64(f), 'f', -1, 32))
	}
	b.WriteByte(']')
	return b.String()
}

var _ VectorDB = (*pgBackend)(nil)
