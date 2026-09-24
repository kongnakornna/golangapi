package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"icmongolang/config"
	"icmongolang/pkg/logger"
)

// ErrNotFound is returned when a document does not exist in the index.
var ErrNotFound = errors.New("elasticsearch: document not found")

// VectorDoc is the document stored in Elasticsearch.
type VectorDoc struct {
	ID         string    `json:"-"`
	DocumentID string    `json:"document_id"`
	Content    string    `json:"content"`
	Embedding  []float32 `json:"embedding"`
	SourceType string    `json:"source_type"`
	Metadata   string    `json:"metadata"`
	Timestamp  time.Time `json:"timestamp,omitempty"`
}

// SearchHit is a single vector-search result.
type SearchHit struct {
	DocumentID string    `json:"document_id"`
	Content    string    `json:"content"`
	SourceType string    `json:"source_type"`
	Metadata   string    `json:"metadata"`
	Timestamp  time.Time `json:"timestamp,omitempty"`
	Score      float64   `json:"score"`
}

// Client is a thin HTTP client for Elasticsearch (REST API).
type Client struct {
	httpClient *http.Client
	addresses  []string
	username   string
	password   string
	index      string
	logger     logger.Logger
}

// NewClient creates an Elasticsearch client and verifies connectivity.
func NewClient(cfg *config.ElasticsearchConfig, log logger.Logger) (*Client, error) {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 10
	}

	c := &Client{
		httpClient: &http.Client{Timeout: time.Duration(timeout) * time.Second},
		addresses:  cfg.Addresses,
		username:   cfg.Username,
		password:   cfg.Password,
		index:      cfg.Index,
		logger:     log,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	if err := c.Ping(ctx); err != nil {
		return nil, fmt.Errorf("elasticsearch connection failed: %w", err)
	}
	return c, nil
}

// WithIndex returns a copy of the client bound to a different index.
// If index is empty it falls back to the default "vector_documents".
func (c *Client) WithIndex(index string) *Client {
	clone := *c
	if index == "" {
		index = "vector_documents"
	}
	clone.index = index
	return &clone
}

// IndexName returns the index name this client is bound to.
func (c *Client) IndexName() string { return c.index }

// Ping checks Elasticsearch availability (GET /).
func (c *Client) Ping(ctx context.Context) error {
	var out map[string]interface{}
	_, err := c.do(ctx, http.MethodGet, "/", nil, &out)
	return err
}

// Health returns Elasticsearch cluster status.
func (c *Client) Health(ctx context.Context) (map[string]interface{}, error) {
	var out map[string]interface{}
	if _, err := c.do(ctx, http.MethodGet, "/_cluster/health", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// EnsureIndex creates the vector index if it does not exist.
// Creating an already-existing index returns 400 resource_already_exists_exception,
// which is treated as success.
func (c *Client) EnsureIndex(ctx context.Context, dims int) error {
	if dims <= 0 {
		dims = 1536
	}

	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"document_id": map[string]interface{}{"type": "keyword"},
				"content":     map[string]interface{}{"type": "text"},
				"embedding": map[string]interface{}{
					"type":       "dense_vector",
					"dims":       dims,
					"index":      true,
					"similarity": "cosine",
				},
				"source_type": map[string]interface{}{"type": "keyword"},
				"metadata":    map[string]interface{}{"type": "text"},
				"timestamp":   map[string]interface{}{"type": "date"},
			},
		},
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(mapping); err != nil {
		return err
	}

	var out map[string]interface{}
	status, err := c.do(ctx, http.MethodPut, "/"+c.index, &body, &out)
	if err != nil {
		// do() returns an error for any status >= 400, so the
		// resource_already_exists_exception case must be matched here.
		if status == http.StatusBadRequest && strings.Contains(err.Error(), "resource_already_exists_exception") {
			c.logger.Infof("✅ Elasticsearch index %q already exists", c.index)
			return nil
		}
		return err
	}

	c.logger.Infof("✅ Elasticsearch index %q ready", c.index)
	return nil
}

// Index stores a vector document and returns the generated document id.
func (c *Client) Index(ctx context.Context, doc *VectorDoc) (string, error) {
	if doc.Timestamp.IsZero() {
		doc.Timestamp = time.Now()
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(doc); err != nil {
		return "", err
	}

	var out struct {
		ID     string `json:"_id"`
		Result string `json:"result"`
	}
	if _, err := c.do(ctx, http.MethodPost, "/"+c.index+"/_doc", &body, &out); err != nil {
		return "", err
	}

	if out.ID == "" {
		return "", fmt.Errorf("elasticsearch returned no document id")
	}
	doc.ID = out.ID
	return out.ID, nil
}

// IndexWithID stores a vector document under an explicit id.
func (c *Client) IndexWithID(ctx context.Context, id string, doc *VectorDoc) error {
	if doc.Timestamp.IsZero() {
		doc.Timestamp = time.Now()
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(doc); err != nil {
		return err
	}

	var out map[string]interface{}
	_, err := c.do(ctx, http.MethodPut, "/"+c.index+"/_doc/"+id, &body, &out)
	return err
}

// SearchKNN performs cosine-similarity vector search over the index.
func (c *Client) SearchKNN(ctx context.Context, embedding []float32, k int) ([]SearchHit, error) {
	if k <= 0 {
		k = 10
	}
	numCandidates := k * 10
	if numCandidates < 100 {
		numCandidates = 100
	}

	query := map[string]interface{}{
		"knn": map[string]interface{}{
			"field":          "embedding",
			"query_vector":   embedding,
			"k":              k,
			"num_candidates": numCandidates,
		},
		"size": k,
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(query); err != nil {
		return nil, err
	}

	var out struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				ID     string          `json:"_id"`
				Score  float64         `json:"_score"`
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if _, err := c.do(ctx, http.MethodPost, "/"+c.index+"/_search", &body, &out); err != nil {
		return nil, err
	}

	hits := make([]SearchHit, 0, len(out.Hits.Hits))
	for _, h := range out.Hits.Hits {
		hit := SearchHit{Score: h.Score}
		// _source is the document itself
		var src VectorDoc
		if err := json.Unmarshal(h.Source, &src); err != nil {
			c.logger.Warnf("⚠️ failed to decode hit %s: %v", h.ID, err)
			continue
		}
		hit.DocumentID = src.DocumentID
		hit.Content = src.Content
		hit.SourceType = src.SourceType
		hit.Metadata = src.Metadata
		hit.Timestamp = src.Timestamp
		hits = append(hits, hit)
	}

	return hits, nil
}

// Delete removes a document by id.
func (c *Client) Delete(ctx context.Context, id string) error {
	var out map[string]interface{}
	_, err := c.do(ctx, http.MethodDelete, "/"+c.index+"/_doc/"+id, nil, &out)
	if err != nil {
		if isNotFoundErr(err) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

// Get returns a document by its id.
func (c *Client) Get(ctx context.Context, id string) (*VectorDoc, error) {
	var out struct {
		Found  bool            `json:"found"`
		Source json.RawMessage `json:"_source"`
	}
	if _, err := c.do(ctx, http.MethodGet, "/"+c.index+"/_doc/"+id, nil, &out); err != nil {
		if isNotFoundErr(err) {
			return nil, nil
		}
		return nil, err
	}
	if !out.Found {
		return nil, nil
	}
	var doc VectorDoc
	if err := json.Unmarshal(out.Source, &doc); err != nil {
		return nil, fmt.Errorf("elasticsearch decode document: %w", err)
	}
	doc.ID = id
	return &doc, nil
}

// List returns a page of documents plus the total count.
func (c *Client) List(ctx context.Context, from, size int) ([]VectorDoc, int64, error) {
	if from < 0 {
		from = 0
	}
	if size <= 0 {
		size = 10
	}

	query := map[string]interface{}{
		"from": from,
		"size": size,
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"sort": []map[string]interface{}{
			{"document_id": "asc"},
		},
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(query); err != nil {
		return nil, 0, err
	}

	var out struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				ID     string          `json:"_id"`
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if _, err := c.do(ctx, http.MethodPost, "/"+c.index+"/_search", &body, &out); err != nil {
		return nil, 0, err
	}

	docs := make([]VectorDoc, 0, len(out.Hits.Hits))
	for _, h := range out.Hits.Hits {
		var doc VectorDoc
		if err := json.Unmarshal(h.Source, &doc); err != nil {
			c.logger.Warnf("⚠️ failed to decode list hit %s: %v", h.ID, err)
			continue
		}
		doc.ID = h.ID
		docs = append(docs, doc)
	}
	return docs, out.Hits.Total.Value, nil
}

// Update merges the non-empty mutable fields of an existing document.
// Empty fields are left untouched (partial update via ES doc merge).
func (c *Client) Update(ctx context.Context, id string, doc *VectorDoc) error {
	if doc == nil {
		return fmt.Errorf("elasticsearch update: doc is nil")
	}

	fields := map[string]interface{}{}
	if doc.DocumentID != "" {
		fields["document_id"] = doc.DocumentID
	}
	if doc.Content != "" {
		fields["content"] = doc.Content
	}
	if len(doc.Embedding) > 0 {
		fields["embedding"] = doc.Embedding
	}
	if doc.SourceType != "" {
		fields["source_type"] = doc.SourceType
	}
	if doc.Metadata != "" {
		fields["metadata"] = doc.Metadata
	}
	if !doc.Timestamp.IsZero() {
		fields["timestamp"] = doc.Timestamp
	}

	if len(fields) == 0 {
		return nil
	}

	update := map[string]interface{}{
		"doc": fields,
	}
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(update); err != nil {
		return err
	}
	var out map[string]interface{}
	_, err := c.do(ctx, http.MethodPost, "/"+c.index+"/_update/"+id, &body, &out)
	if err != nil {
		if isNotFoundErr(err) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

// isNotFoundErr reports whether the error is an HTTP 404 from the ES API.
func isNotFoundErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "returned 404")
}

// do executes an HTTP request against the first reachable address.
func (c *Client) do(ctx context.Context, method, path string, body io.Reader, out interface{}) (int, error) {
	var lastErr error
	for _, addr := range c.addresses {
		url := strings.TrimSuffix(addr, "/") + path

		req, err := http.NewRequestWithContext(ctx, method, url, body)
		if err != nil {
			return 0, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		if c.username != "" || c.password != "" {
			req.SetBasicAuth(c.username, c.password)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode >= http.StatusInternalServerError && len(c.addresses) > 1 {
			lastErr = fmt.Errorf("elasticsearch %s returned %d: %s", path, resp.StatusCode, string(respBody))
			continue
		}

		if out != nil && len(respBody) > 0 {
			if err := json.Unmarshal(respBody, out); err != nil {
				return resp.StatusCode, fmt.Errorf("elasticsearch response decode: %w", err)
			}
		}

		if resp.StatusCode >= 400 {
			return resp.StatusCode, fmt.Errorf("elasticsearch %s returned %d: %s", path, resp.StatusCode, string(respBody))
		}

		return resp.StatusCode, nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("no elasticsearch addresses configured")
	}
	return 0, lastErr
}
