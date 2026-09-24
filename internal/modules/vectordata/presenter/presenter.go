package presenter

// CreateDocumentRequest – create a vector document.
// Provide Content (embedded via LLM) or a pre-computed Embedding.
type CreateDocumentRequest struct {
	DocumentID string    `json:"document_id"`
	Content    string    `json:"content"`
	Embedding  []float32 `json:"embedding"`
	SourceType string    `json:"source_type"`
	Metadata   string    `json:"metadata"`
}

// UpdateDocumentRequest – update an existing vector document.
type UpdateDocumentRequest struct {
	DocumentID string    `json:"document_id"`
	Content    string    `json:"content"`
	Embedding  []float32 `json:"embedding"`
	SourceType string    `json:"source_type"`
	Metadata   string    `json:"metadata"`
}

// DocumentResponse – a single vector document.
type DocumentResponse struct {
	ID         string `json:"id"`
	DocumentID string `json:"document_id"`
	Content    string `json:"content"`
	SourceType string `json:"source_type"`
	Metadata   string `json:"metadata"`
	Dim        int    `json:"dim"`
	Indexed    bool   `json:"indexed"`
}

// ListResponse – paginated list of documents.
type ListResponse struct {
	Total     int64              `json:"total"`
	Limit     int                `json:"limit"`
	Offset    int                `json:"offset"`
	Documents []DocumentResponse `json:"documents"`
}

// SearchRequest – semantic (vector) search query.
type SearchRequest struct {
	Query string `json:"query" validate:"required"`
	K     int    `json:"k"`
}

// SearchHit – single vector-search result.
type SearchHit struct {
	DocumentID string  `json:"document_id"`
	Content    string  `json:"content"`
	SourceType string  `json:"source_type"`
	Metadata   string  `json:"metadata"`
	Score      float64 `json:"score"`
}

// SearchResponse – vector-search result set.
type SearchResponse struct {
	Query string      `json:"query"`
	K     int         `json:"k"`
	Count int         `json:"count"`
	Hits  []SearchHit `json:"hits"`
}

// CreateIndexRequest – create/ensure the vector index with given dimensions.
type CreateIndexRequest struct {
	Dims int `json:"dims"`
}

// SeedRequest – seed sample documents.
type SeedRequest struct {
	Count int `json:"count"`
}

// SeedResponse – result of seeding.
type SeedResponse struct {
	Provider string `json:"provider"`
	Index    string `json:"index"`
	Seeded   int    `json:"seeded"`
}
