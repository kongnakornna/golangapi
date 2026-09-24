package usecase

import (
	"context"
	"fmt"
	"strings"

	"icmongolang/config"
	"icmongolang/internal/modules/vectordata/presenter"
	"icmongolang/pkg/llm"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/vectordb"
)

// VectorDataUseCase orchestrates LLM embedding + abstract vector database ops.
type VectorDataUseCase interface {
	Health(ctx context.Context) (map[string]interface{}, error)
	CreateIndex(ctx context.Context, req *presenter.CreateIndexRequest) error
	Create(ctx context.Context, req *presenter.CreateDocumentRequest) (*presenter.DocumentResponse, error)
	Get(ctx context.Context, id string) (*presenter.DocumentResponse, error)
	List(ctx context.Context, limit, offset int) (*presenter.ListResponse, error)
	Update(ctx context.Context, id string, req *presenter.UpdateDocumentRequest) (*presenter.DocumentResponse, error)
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, req *presenter.SearchRequest) (*presenter.SearchResponse, error)
	Seed(ctx context.Context, req *presenter.SeedRequest) (*presenter.SeedResponse, error)
}

type vectorDataUseCase struct {
	vdb    vectordb.VectorDB
	llm    *llm.Client
	cfg    *config.Config
	logger logger.Logger
	index  string
}

// NewVectorDataUseCase creates the usecase.
func NewVectorDataUseCase(vdb vectordb.VectorDB, llmClient *llm.Client, cfg *config.Config, log logger.Logger) VectorDataUseCase {
	return &vectorDataUseCase{vdb: vdb, llm: llmClient, cfg: cfg, logger: log, index: cfg.VectorDB.Index}
}

func (u *vectorDataUseCase) Health(ctx context.Context) (map[string]interface{}, error) {
	status := map[string]interface{}{
		"provider": u.vdb.Provider(),
		"index":    u.index,
		"status":   "up",
	}
	if err := u.vdb.Ping(ctx); err != nil {
		status["status"] = "down"
		status["error"] = err.Error()
	}
	return status, nil
}

func (u *vectorDataUseCase) CreateIndex(ctx context.Context, req *presenter.CreateIndexRequest) error {
	dims := u.dims()
	if req != nil && req.Dims > 0 {
		dims = req.Dims
	}
	if err := u.vdb.EnsureIndex(ctx, dims); err != nil {
		return err
	}
	u.logger.Infof("✅ Vector index %q ready (provider=%s, dims=%d)", u.index, u.vdb.Provider(), dims)
	return nil
}

// Create stores a document. If Embedding is empty it is generated via the LLM.
func (u *vectorDataUseCase) Create(ctx context.Context, req *presenter.CreateDocumentRequest) (*presenter.DocumentResponse, error) {
	embedding := req.Embedding
	if len(embedding) == 0 {
		var err error
		embedding, err = u.embed(ctx, req.Content)
		if err != nil {
			return nil, err
		}
	}

	if err := u.vdb.EnsureIndex(ctx, len(embedding)); err != nil {
		return nil, fmt.Errorf("ensure index: %w", err)
	}

	doc := &vectordb.VectorDoc{
		DocumentID: req.DocumentID,
		Content:    req.Content,
		Embedding:  embedding,
		SourceType: req.SourceType,
		Metadata:   req.Metadata,
	}
	id, err := u.vdb.Index(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("index vector: %w", err)
	}

	u.logger.Infof("✅ Vector document %q indexed (id=%s, provider=%s, dims=%d)", req.DocumentID, id, u.vdb.Provider(), len(embedding))

	return &presenter.DocumentResponse{
		ID:         id,
		DocumentID: req.DocumentID,
		Content:    req.Content,
		SourceType: req.SourceType,
		Metadata:   req.Metadata,
		Dim:        len(embedding),
		Indexed:    true,
	}, nil
}

func (u *vectorDataUseCase) Get(ctx context.Context, id string) (*presenter.DocumentResponse, error) {
	doc, err := u.vdb.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &presenter.DocumentResponse{
		ID:         doc.ID,
		DocumentID: doc.DocumentID,
		Content:    doc.Content,
		SourceType: doc.SourceType,
		Metadata:   doc.Metadata,
	}, nil
}

func (u *vectorDataUseCase) List(ctx context.Context, limit, offset int) (*presenter.ListResponse, error) {
	docs, total, err := u.vdb.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	resp := &presenter.ListResponse{Total: total, Limit: limit, Offset: offset, Documents: make([]presenter.DocumentResponse, 0, len(docs))}
	for _, d := range docs {
		resp.Documents = append(resp.Documents, presenter.DocumentResponse{
			ID:         d.ID,
			DocumentID: d.DocumentID,
			Content:    d.Content,
			SourceType: d.SourceType,
			Metadata:   d.Metadata,
		})
	}
	return resp, nil
}

// Update partially updates an existing document. Only non-empty request
// fields are applied; the rest are preserved from the stored document.
// If Content is provided and Embedding is not, a new embedding is generated.
func (u *vectorDataUseCase) Update(ctx context.Context, id string, req *presenter.UpdateDocumentRequest) (*presenter.DocumentResponse, error) {
	existing, err := u.vdb.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	embedding := req.Embedding
	content := req.Content
	documentID := req.DocumentID
	sourceType := req.SourceType
	metadata := req.Metadata

	if len(embedding) == 0 && strings.TrimSpace(content) != "" {
		embedding, err = u.embed(ctx, content)
		if err != nil {
			return nil, err
		}
	}

	// Merge: only override fields that were actually provided.
	if content == "" {
		content = existing.Content
	}
	if documentID == "" {
		documentID = existing.DocumentID
	}
	if sourceType == "" {
		sourceType = existing.SourceType
	}
	if metadata == "" {
		metadata = existing.Metadata
	}

	doc := &vectordb.VectorDoc{
		DocumentID: documentID,
		Content:    content,
		Embedding:  embedding,
		SourceType: sourceType,
		Metadata:   metadata,
	}
	if err := u.vdb.Update(ctx, id, doc); err != nil {
		return nil, err
	}

	u.logger.Infof("✅ Vector document %s updated (provider=%s)", id, u.vdb.Provider())

	return &presenter.DocumentResponse{
		ID:         id,
		DocumentID: documentID,
		Content:    content,
		SourceType: sourceType,
		Metadata:   metadata,
		Dim:        len(embedding),
		Indexed:    true,
	}, nil
}

func (u *vectorDataUseCase) Delete(ctx context.Context, id string) error {
	if err := u.vdb.Delete(ctx, id); err != nil {
		return err
	}
	u.logger.Infof("🗑️ Vector document %s deleted (provider=%s)", id, u.vdb.Provider())
	return nil
}

// Search embeds the query and runs cosine-similarity search.
func (u *vectorDataUseCase) Search(ctx context.Context, req *presenter.SearchRequest) (*presenter.SearchResponse, error) {
	if u.llm == nil {
		return nil, fmt.Errorf("llm client is nil – cannot embed query")
	}
	k := req.K
	if k <= 0 {
		k = 10
	}
	if k > 100 {
		k = 100
	}

	embedding, err := u.llm.Embed(ctx, req.Query)
	if err != nil {
		return nil, fmt.Errorf("embed query: %w", err)
	}

	hits, err := u.vdb.SearchKNN(ctx, embedding, k)
	if err != nil {
		return nil, fmt.Errorf("vector search: %w", err)
	}

	resp := &presenter.SearchResponse{Query: req.Query, K: k, Count: len(hits), Hits: make([]presenter.SearchHit, 0, len(hits))}
	for _, h := range hits {
		resp.Hits = append(resp.Hits, presenter.SearchHit{
			DocumentID: h.DocumentID,
			Content:    h.Content,
			SourceType: h.SourceType,
			Metadata:   h.Metadata,
			Score:      h.Score,
		})
	}

	u.logger.Infof("🔍 Vector search %q returned %d hits (provider=%s)", req.Query, len(hits), u.vdb.Provider())
	return resp, nil
}

// Seed indexes a set of sample documents through the LLM. It is idempotent:
// samples whose document_id already exists are skipped.
func (u *vectorDataUseCase) Seed(ctx context.Context, req *presenter.SeedRequest) (*presenter.SeedResponse, error) {
	count := len(sampleDocs)
	if req != nil && req.Count > 0 {
		count = req.Count
	}
	if count > len(sampleDocs) {
		count = len(sampleDocs)
	}

	existing, _, err := u.vdb.List(ctx, 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("seed: list existing documents: %w", err)
	}
	seen := make(map[string]bool, len(existing))
	for _, d := range existing {
		if d.DocumentID != "" {
			seen[d.DocumentID] = true
		}
	}

	seeded := 0
	for i := 0; i < count; i++ {
		s := sampleDocs[i]
		if seen[s.DocumentID] {
			u.logger.Infof("🌱 Skipping sample %q (already exists)", s.DocumentID)
			continue
		}
		embedding, err := u.embed(ctx, s.Content)
		if err != nil {
			return nil, fmt.Errorf("embed sample %d: %w", i, err)
		}
		if err := u.vdb.EnsureIndex(ctx, len(embedding)); err != nil {
			return nil, fmt.Errorf("ensure index: %w", err)
		}
		if _, err := u.vdb.Index(ctx, &vectordb.VectorDoc{
			DocumentID: s.DocumentID,
			Content:    s.Content,
			Embedding:  embedding,
			SourceType: s.SourceType,
		}); err != nil {
			return nil, fmt.Errorf("index sample %d: %w", i, err)
		}
		seeded++
		seen[s.DocumentID] = true
	}

	u.logger.Infof("🌱 Seeded %d sample documents into %q (provider=%s)", seeded, u.index, u.vdb.Provider())
	return &presenter.SeedResponse{Provider: u.vdb.Provider(), Index: u.index, Seeded: seeded}, nil
}

// embed generates an embedding via the configured LLM.
func (u *vectorDataUseCase) embed(ctx context.Context, text string) ([]float32, error) {
	if u.llm == nil {
		return nil, fmt.Errorf("llm client is nil – cannot generate embedding")
	}
	embedding, err := u.llm.Embed(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("embed content: %w", err)
	}
	return embedding, nil
}

func (u *vectorDataUseCase) dims() int {
	if u.cfg.VectorDB.Dims > 0 {
		return u.cfg.VectorDB.Dims
	}
	if u.llm != nil {
		return u.llm.Dims()
	}
	return 768
}

// sampleDocs – example IoT/sensor knowledge-base documents for seeding.
var sampleDocs = []struct {
	DocumentID string
	SourceType string
	Content    string
}{
	{
		DocumentID: "doc-temp-sensor",
		SourceType: "iot-manual",
		Content:    "เซนเซอร์วัดอุณหภูมิ (temperature sensor) ใช้สำหรับตรวจวัดอุณหภูมิอากาศหรือน้ำในระบบ IoT ถ้าค่าเกิน 60 องศาเซนติเกรด จะส่งสัญญาณแจ้งเตือนไปยังระบบ Worker ผ่าน MQTT",
	},
	{
		DocumentID: "doc-humidity-sensor",
		SourceType: "iot-manual",
		Content:    "เซนเซอร์วัดความชื้น (humidity sensor) วัดความชื้นสัมพัทธ์ของอากาศ ถ้าความชื้นต่ำกว่า 30% จะแจ้งเตือนให้เปิดระบบพ่นหมอกเพื่อป้องกันฝุ่น",
	},
	{
		DocumentID: "doc-vibration-sensor",
		SourceType: "iot-manual",
		Content:    "เซนเซอร์วัดแรงสั่นสะเทือน (vibration sensor) ใช้ตรวจสอบสุขภาพของมอเตอร์และเครื่องจักร ถ้าค่าแรงสั่นสะเทือนสูงผิดปกติ ให้แจ้งวิศวกรตรวจสอบทันที",
	},
	{
		DocumentID: "doc-power-meter",
		SourceType: "iot-manual",
		Content:    "เครื่องวัดพลังงานไฟฟ้า (power meter) วัดค่าแรงดัน กระแส และกำลังไฟฟ้าของเครื่องจักร เพื่อใช้วิเคราะห์ประสิทธิภาพพลังงานและค่าไฟรายเดือน",
	},
	{
		DocumentID: "doc-alarm-policy",
		SourceType: "ops-policy",
		Content:    "นโยบายการแจ้งเตือน: ระดับ Critical ต้องแจ้งภายใน 1 นาที, ระดับ Warning แจ้งภายใน 5 นาที, ทุกระดับต้องบันทึกลง alarm log และส่งผ่าน WebSocket ไปยัง Dashboard",
	},
}
