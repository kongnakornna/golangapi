package usecase

import (
	"context"
	"fmt"

	"icmongolang/internal/modules/elasticsearch/presenter"
	esclient "icmongolang/pkg/elasticsearch"
	"icmongolang/pkg/llm"
	"icmongolang/pkg/logger"
)

// ElasticsearchUseCase orchestrates LLM embedding + Elasticsearch vector search.
type ElasticsearchUseCase interface {
	Health(ctx context.Context) (map[string]interface{}, error)
	CreateIndex(ctx context.Context, req *presenter.CreateIndexRequest) error
	EmbedAndIndex(ctx context.Context, req *presenter.EmbeddingRequest) (*presenter.EmbeddingResponse, error)
	IndexVector(ctx context.Context, req *presenter.IndexVectorRequest) (*presenter.EmbeddingResponse, error)
	SearchVector(ctx context.Context, req *presenter.SearchRequest) (*presenter.SearchResponse, error)
}

type elasticsearchUseCase struct {
	es     *esclient.Client
	llm    *llm.Client
	logger logger.Logger
}

// NewElasticsearchUseCase creates the usecase.
func NewElasticsearchUseCase(es *esclient.Client, llmClient *llm.Client, log logger.Logger) ElasticsearchUseCase {
	return &elasticsearchUseCase{es: es, llm: llmClient, logger: log}
}

func (u *elasticsearchUseCase) Health(ctx context.Context) (map[string]interface{}, error) {
	return u.es.Health(ctx)
}

func (u *elasticsearchUseCase) CreateIndex(ctx context.Context, req *presenter.CreateIndexRequest) error {
	if u.llm == nil {
		return fmt.Errorf("llm client is nil – cannot determine index dims")
	}
	dims := u.llm.Dims()
	if req != nil && req.Dims > 0 {
		dims = req.Dims
	}
	if err := u.es.EnsureIndex(ctx, dims); err != nil {
		return err
	}
	u.logger.Infof("✅ Vector index %q created (dims=%d)", u.es.IndexName(), dims)
	return nil
}

// EmbedAndIndex generates an embedding via the LLM and stores it in Elasticsearch.
func (u *elasticsearchUseCase) EmbedAndIndex(ctx context.Context, req *presenter.EmbeddingRequest) (*presenter.EmbeddingResponse, error) {
	if u.llm == nil {
		return nil, fmt.Errorf("llm client is nil – cannot generate embedding")
	}

	embedding, err := u.llm.Embed(ctx, req.Content)
	if err != nil {
		return nil, fmt.Errorf("embed content: %w", err)
	}

	if err := u.es.EnsureIndex(ctx, len(embedding)); err != nil {
		return nil, fmt.Errorf("ensure index: %w", err)
	}

	doc := &esclient.VectorDoc{
		DocumentID: req.DocumentID,
		Content:    req.Content,
		Embedding:  embedding,
		SourceType: req.SourceType,
		Metadata:   req.Metadata,
	}

	id, err := u.es.Index(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("index vector: %w", err)
	}

	u.logger.Infof("✅ Embedded & indexed document %q (id=%s, dims=%d)", req.DocumentID, id, len(embedding))

	return &presenter.EmbeddingResponse{
		DocumentID: req.DocumentID,
		Content:    req.Content,
		SourceType: req.SourceType,
		Index:      u.es.IndexName(),
		Model:      u.llm.Model(),
		Dim:        len(embedding),
		Indexed:    true,
	}, nil
}

// IndexVector stores a pre-computed vector directly.
func (u *elasticsearchUseCase) IndexVector(ctx context.Context, req *presenter.IndexVectorRequest) (*presenter.EmbeddingResponse, error) {
	if err := u.es.EnsureIndex(ctx, len(req.Embedding)); err != nil {
		return nil, fmt.Errorf("ensure index: %w", err)
	}

	doc := &esclient.VectorDoc{
		DocumentID: req.DocumentID,
		Content:    req.Content,
		Embedding:  req.Embedding,
		SourceType: req.SourceType,
		Metadata:   req.Metadata,
	}

	id, err := u.es.Index(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("index vector: %w", err)
	}

	u.logger.Infof("✅ Indexed vector document %q (id=%s, dims=%d)", req.DocumentID, id, len(req.Embedding))

	return &presenter.EmbeddingResponse{
		DocumentID: req.DocumentID,
		Content:    req.Content,
		SourceType: req.SourceType,
		Index:      u.es.IndexName(),
		Dim:        len(req.Embedding),
		Indexed:    true,
	}, nil
}

// SearchVector embeds the query and runs cosine-similarity search.
func (u *elasticsearchUseCase) SearchVector(ctx context.Context, req *presenter.SearchRequest) (*presenter.SearchResponse, error) {
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

	hits, err := u.es.SearchKNN(ctx, embedding, k)
	if err != nil {
		return nil, fmt.Errorf("vector search: %w", err)
	}

	resp := &presenter.SearchResponse{Query: req.Query, K: k}
	for _, h := range hits {
		resp.Hits = append(resp.Hits, presenter.SearchHit{
			DocumentID: h.DocumentID,
			Content:    h.Content,
			SourceType: h.SourceType,
			Metadata:   h.Metadata,
			Score:      h.Score,
		})
	}

	u.logger.Infof("🔍 Vector search %q returned %d hits", req.Query, len(hits))
	return resp, nil
}
