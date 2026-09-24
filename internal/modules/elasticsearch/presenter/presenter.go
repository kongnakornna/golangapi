package presenter

// EmbeddingRequest – content to embed via the LLM and store in Elasticsearch.
type EmbeddingRequest struct {
	DocumentID string `json:"document_id"`
	Content    string `json:"content"`
	SourceType string `json:"source_type"`
	Metadata   string `json:"metadata"`
}

// IndexVectorRequest – already-computed vector to store directly.
type IndexVectorRequest struct {
	DocumentID string    `json:"document_id"`
	Content    string    `json:"content"`
	Embedding  []float32 `json:"embedding"`
	SourceType string    `json:"source_type"`
	Metadata   string    `json:"metadata"`
}

// EmbeddingResponse – result of embedding + indexing.
type EmbeddingResponse struct {
	DocumentID string `json:"document_id"`
	Content    string `json:"content"`
	SourceType string `json:"source_type"`
	Index      string `json:"index"`
	Model      string `json:"model,omitempty"`
	Dim        int    `json:"dim"`
	Indexed    bool   `json:"indexed"`
}

// SearchRequest – semantic (vector) search query.
type SearchRequest struct {
	Query string `json:"query"`
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
	Hits  []SearchHit `json:"hits"`
}

// CreateIndexRequest – create the vector index with given dimensions.
type CreateIndexRequest struct {
	Dims int `json:"dims"`
}
