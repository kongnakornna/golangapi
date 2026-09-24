package vectordb

import (
	"context"
	"errors"
	"time"

	"icmongolang/config"
	esclient "icmongolang/pkg/elasticsearch"
	"icmongolang/pkg/logger"
)

// esBackend implements VectorDB on top of the Elasticsearch dense_vector index.
type esBackend struct {
	client *esclient.Client
	index  string
	log    logger.Logger
}

func newESBackend(cfg *config.VectorDBConfig, es *esclient.Client, log logger.Logger) VectorDB {
	index := cfg.Index
	if index == "" {
		index = "vector_documents"
	}
	// The shared esclient is bound to cfg.Elasticsearch.Index, so bind a
	// per-backend client to the vector DB index (BLOCKER-1 fix).
	return &esBackend{client: es.WithIndex(index), index: index, log: log}
}

func (b *esBackend) Provider() string { return "elasticsearch" }

func (b *esBackend) Ping(ctx context.Context) error {
	return b.client.Ping(ctx)
}

func (b *esBackend) EnsureIndex(ctx context.Context, dims int) error {
	if dims <= 0 {
		dims = 768
	}
	return b.client.EnsureIndex(ctx, dims)
}

func (b *esBackend) Index(ctx context.Context, doc *VectorDoc) (string, error) {
	esDoc := &esclient.VectorDoc{
		DocumentID: doc.DocumentID,
		Content:    doc.Content,
		Embedding:  doc.Embedding,
		SourceType: doc.SourceType,
		Metadata:   doc.Metadata,
		Timestamp:  doc.Timestamp,
	}
	if esDoc.Timestamp.IsZero() {
		esDoc.Timestamp = time.Now()
	}
	id, err := b.client.Index(ctx, esDoc)
	if err != nil {
		return "", err
	}
	doc.ID = id
	doc.Timestamp = esDoc.Timestamp
	return id, nil
}

func (b *esBackend) Get(ctx context.Context, id string) (*VectorDoc, error) {
	esDoc, err := b.client.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if esDoc == nil {
		return nil, ErrNotFound
	}
	return &VectorDoc{
		ID:         esDoc.ID,
		DocumentID: esDoc.DocumentID,
		Content:    esDoc.Content,
		Embedding:  esDoc.Embedding,
		SourceType: esDoc.SourceType,
		Metadata:   esDoc.Metadata,
		Timestamp:  esDoc.Timestamp,
	}, nil
}

func (b *esBackend) List(ctx context.Context, limit, offset int) ([]VectorDoc, int64, error) {
	esDocs, total, err := b.client.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	docs := make([]VectorDoc, 0, len(esDocs))
	for _, d := range esDocs {
		docs = append(docs, VectorDoc{
			ID:         d.ID,
			DocumentID: d.DocumentID,
			Content:    d.Content,
			Embedding:  d.Embedding,
			SourceType: d.SourceType,
			Metadata:   d.Metadata,
			Timestamp:  d.Timestamp,
		})
	}
	return docs, total, nil
}

func (b *esBackend) Update(ctx context.Context, id string, doc *VectorDoc) error {
	esDoc := &esclient.VectorDoc{
		DocumentID: doc.DocumentID,
		Content:    doc.Content,
		Embedding:  doc.Embedding,
		SourceType: doc.SourceType,
		Metadata:   doc.Metadata,
		Timestamp:  doc.Timestamp,
	}
	err := b.client.Update(ctx, id, esDoc)
	if errors.Is(err, esclient.ErrNotFound) {
		return ErrNotFound
	}
	return err
}

func (b *esBackend) Delete(ctx context.Context, id string) error {
	err := b.client.Delete(ctx, id)
	if errors.Is(err, esclient.ErrNotFound) {
		return ErrNotFound
	}
	return err
}

func (b *esBackend) SearchKNN(ctx context.Context, embedding []float32, k int) ([]SearchHit, error) {
	hits, err := b.client.SearchKNN(ctx, embedding, k)
	if err != nil {
		return nil, err
	}
	out := make([]SearchHit, 0, len(hits))
	for _, h := range hits {
		out = append(out, SearchHit{
			DocumentID: h.DocumentID,
			Content:    h.Content,
			SourceType: h.SourceType,
			Metadata:   h.Metadata,
			Score:      h.Score,
			Timestamp:  h.Timestamp,
		})
	}
	return out, nil
}

var _ VectorDB = (*esBackend)(nil)
