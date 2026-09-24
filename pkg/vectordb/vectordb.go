// Package vectordb provides an abstraction over vector databases.
//
// Supported providers:
//   - "elasticsearch": dense_vector index in Elasticsearch
//   - "pgvector":      PostgreSQL + pgvector extension
package vectordb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"icmongolang/config"
	esclient "icmongolang/pkg/elasticsearch"
	"icmongolang/pkg/logger"

	"gorm.io/gorm"
)

// ErrNotFound is returned when a document does not exist in the backend.
var ErrNotFound = errors.New("vectordb: document not found")

// VectorDoc is a document stored in the vector database.
type VectorDoc struct {
	ID         string    `json:"id"`
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
	Score      float64   `json:"score"`
	Timestamp  time.Time `json:"timestamp,omitempty"`
}

// VectorDB is the backend-agnostic interface used by the vectordata module.
type VectorDB interface {
	// Provider returns the backend name ("elasticsearch" | "pgvector").
	Provider() string
	// Ping verifies connectivity to the backend.
	Ping(ctx context.Context) error
	// EnsureIndex creates the index/table if it does not exist.
	EnsureIndex(ctx context.Context, dims int) error
	// Index stores a document (embedding required) and returns its id.
	Index(ctx context.Context, doc *VectorDoc) (string, error)
	// Get returns a document by its id.
	Get(ctx context.Context, id string) (*VectorDoc, error)
	// List returns a page of documents plus the total count.
	List(ctx context.Context, limit, offset int) ([]VectorDoc, int64, error)
	// Update replaces the fields of an existing document.
	Update(ctx context.Context, id string, doc *VectorDoc) error
	// Delete removes a document by id.
	Delete(ctx context.Context, id string) error
	// SearchKNN performs cosine-similarity vector search.
	SearchKNN(ctx context.Context, embedding []float32, k int) ([]SearchHit, error)
}

// New builds a VectorDB backend from config.
//
//	cfg.Provider == "elasticsearch" → Elasticsearch backend (requires es != nil)
//	cfg.Provider == "pgvector"     → PostgreSQL + pgvector backend (requires db != nil)
func New(cfg *config.VectorDBConfig, es *esclient.Client, db *gorm.DB, log logger.Logger) (VectorDB, error) {
	switch strings.ToLower(cfg.Provider) {
	case "elasticsearch":
		if es == nil {
			return nil, fmt.Errorf("vectordb: provider elasticsearch but es client is nil")
		}
		return newESBackend(cfg, es, log), nil
	case "pgvector", "":
		if db == nil {
			return nil, fmt.Errorf("vectordb: provider pgvector but db is nil")
		}
		return newPgBackend(cfg, db, log), nil
	default:
		return nil, fmt.Errorf("vectordb: unknown provider %q", cfg.Provider)
	}
}
