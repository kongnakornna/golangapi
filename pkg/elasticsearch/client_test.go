package elasticsearch

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"icmongolang/config"

	"github.com/stretchr/testify/assert"
)

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func mustNotNil(t *testing.T, v interface{}) {
	t.Helper()
	if v == nil {
		t.Fatal("expected non-nil")
	}
}

type dummyLogger struct{}

func (dummyLogger) InitLogger()                        {}
func (dummyLogger) Debug(args ...interface{})          {}
func (dummyLogger) Debugf(t string, a ...interface{})  {}
func (dummyLogger) Info(args ...interface{})           {}
func (dummyLogger) Infof(t string, a ...interface{})   {}
func (dummyLogger) Warn(args ...interface{})           {}
func (dummyLogger) Warnf(t string, a ...interface{})   {}
func (dummyLogger) Error(args ...interface{})          {}
func (dummyLogger) Errorf(t string, a ...interface{})  {}
func (dummyLogger) DPanic(args ...interface{})         {}
func (dummyLogger) DPanicf(t string, a ...interface{}) {}
func (dummyLogger) Fatal(args ...interface{})          {}
func (dummyLogger) Fatalf(t string, a ...interface{})  {}
func (dummyLogger) Sync() error                        { return nil }

type recorder struct {
	updateBodies []string
	lastIndex    string
}

func newFakeServer(t *testing.T, rec *recorder) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/":
			_, _ = w.Write([]byte(`{"tagline":"You Know, for Search"}`))
		case r.Method == http.MethodPut && strings.Count(strings.Trim(r.URL.Path, "/"), "/") == 0:
			rec.lastIndex = strings.Trim(r.URL.Path, "/")
			_, _ = w.Write([]byte(`{"acknowledged":true}`))
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/_doc"):
			_, _ = w.Write([]byte(`{"_id":"doc-1","result":"created"}`))
		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/_update/"):
			body, _ := io.ReadAll(r.Body)
			rec.updateBodies = append(rec.updateBodies, string(body))
			if strings.HasSuffix(r.URL.Path, "/missing") {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"error":{"type":"document_missing_exception"}}`))
				return
			}
			_, _ = w.Write([]byte(`{"result":"updated"}`))
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/_search"):
			_, _ = w.Write([]byte(`{"hits":{"total":{"value":2},"hits":[
				{"_id":"doc-1","_score":0.9,"_source":{"document_id":"a","content":"hello","source_type":"t","metadata":"m","timestamp":"2026-01-01T00:00:00Z"}},
				{"_id":"doc-2","_score":0.5,"_source":{"document_id":"b","content":"world","source_type":"t","metadata":"","timestamp":"2026-01-01T00:00:00Z"}}
			]}}`))
		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/_doc/"):
			id := strings.Split(r.URL.Path, "/_doc/")[1]
			if id == "missing" {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"found":false}`))
				return
			}
			_, _ = w.Write([]byte(`{"found":true,"_source":{"document_id":"a","content":"hello","embedding":[1,2,3],"source_type":"t","metadata":"m"}}`))
		case r.Method == http.MethodDelete && strings.Contains(r.URL.Path, "/_doc/"):
			if strings.HasSuffix(r.URL.Path, "/missing") {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"result":"not_found"}`))
				return
			}
			_, _ = w.Write([]byte(`{"result":"deleted"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"not found"}`))
		}
	})
	return httptest.NewServer(mux)
}

func newTestClient(t *testing.T, srv *httptest.Server) *Client {
	t.Helper()
	c, err := NewClient(&config.ElasticsearchConfig{
		Addresses: []string{srv.URL},
		Index:     "vector_embeddings",
		Timeout:   5,
	}, dummyLogger{})
	must(t, err)
	return c
}

func TestClient_Index_Get_RoundTrip(t *testing.T) {
	rec := &recorder{}
	srv := newFakeServer(t, rec)
	defer srv.Close()

	c := newTestClient(t, srv)
	ctx := context.Background()

	doc := &VectorDoc{
		DocumentID: "a",
		Content:    "hello",
		Embedding:  []float32{1, 2, 3},
		SourceType: "t",
		Metadata:   "m",
	}
	id, err := c.Index(ctx, doc)
	must(t, err)
	assert.Equal(t, "doc-1", id)
	assert.Equal(t, "doc-1", doc.ID)

	got, err := c.Get(ctx, id)
	must(t, err)
	mustNotNil(t, got)
	assert.Equal(t, "hello", got.Content)
	assert.Equal(t, "m", got.Metadata)
	assert.Equal(t, id, got.ID)
}

func TestClient_Get_NotFound_ReturnsNil(t *testing.T) {
	rec := &recorder{}
	srv := newFakeServer(t, rec)
	defer srv.Close()

	c := newTestClient(t, srv)
	got, err := c.Get(context.Background(), "missing")
	must(t, err)
	assert.Nil(t, got)
}

func TestClient_Update_PartialOnly(t *testing.T) {
	rec := &recorder{}
	srv := newFakeServer(t, rec)
	defer srv.Close()

	c := newTestClient(t, srv)
	err := c.Update(context.Background(), "doc-1", &VectorDoc{Metadata: "m2"})
	must(t, err)
	if len(rec.updateBodies) != 1 {
		t.Fatalf("expected len(rec.updateBodies)==1, got %d", len(rec.updateBodies))
	}

	var body map[string]map[string]interface{}
	must(t, json.Unmarshal([]byte(rec.updateBodies[0]), &body))
	doc := body["doc"]
	// Only the provided field must be sent; empty fields must be omitted.
	assert.Equal(t, "m2", doc["metadata"])
	assert.NotContains(t, doc, "content")
	assert.NotContains(t, doc, "embedding")
	assert.NotContains(t, doc, "document_id")
	assert.NotContains(t, doc, "source_type")
}

func TestClient_Update_NotFound_ReturnsErrNotFound(t *testing.T) {
	rec := &recorder{}
	srv := newFakeServer(t, rec)
	defer srv.Close()

	c := newTestClient(t, srv)
	err := c.Update(context.Background(), "missing", &VectorDoc{Content: "x"})
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestClient_Delete_NotFound_ReturnsErrNotFound(t *testing.T) {
	rec := &recorder{}
	srv := newFakeServer(t, rec)
	defer srv.Close()

	c := newTestClient(t, srv)
	err := c.Delete(context.Background(), "missing")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestClient_List_And_Search_Metadata(t *testing.T) {
	rec := &recorder{}
	srv := newFakeServer(t, rec)
	defer srv.Close()

	c := newTestClient(t, srv)
	ctx := context.Background()

	docs, total, err := c.List(ctx, 0, 10)
	must(t, err)
	assert.Equal(t, int64(2), total)
	if len(docs) != 2 {
		t.Fatalf("expected len(docs)==2, got %d", len(docs))
	}
	assert.Equal(t, "m", docs[0].Metadata)
	assert.Equal(t, "doc-1", docs[0].ID)

	hits, err := c.SearchKNN(ctx, []float32{1, 2, 3}, 5)
	must(t, err)
	if len(hits) != 2 {
		t.Fatalf("expected len(hits)==2, got %d", len(hits))
	}
	assert.Equal(t, 0.9, hits[0].Score)
	assert.Equal(t, "m", hits[0].Metadata)
}

func TestClient_WithIndex_Isolation(t *testing.T) {
	rec := &recorder{}
	srv := newFakeServer(t, rec)
	defer srv.Close()

	c := newTestClient(t, srv)
	other := c.WithIndex("vector_documents")

	ctx := context.Background()
	must(t, other.EnsureIndex(ctx, 768))
	assert.Equal(t, "vector_documents", rec.lastIndex)

	// EnsureIndex must default the index when empty.
	other2 := c.WithIndex("")
	must(t, other2.EnsureIndex(ctx, 768))
	assert.Equal(t, "vector_documents", rec.lastIndex)
}

func TestClient_EnsureIndex_AlreadyExists(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPut {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":{"type":"resource_already_exists_exception"}}`))
			return
		}
		_, _ = w.Write([]byte(`{}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := newTestClient(t, srv)
	err := c.EnsureIndex(context.Background(), 768)
	assert.NoError(t, err)
}
