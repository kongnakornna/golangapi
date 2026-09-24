package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"icmongolang/config"
	"icmongolang/internal/modules/vectordata/presenter"
	"icmongolang/pkg/llm"
	"icmongolang/pkg/vectordb"

	"github.com/stretchr/testify/assert"
)

type noopLogger struct{}

func (noopLogger) InitLogger()                        {}
func (noopLogger) Debug(args ...interface{})          {}
func (noopLogger) Debugf(t string, a ...interface{})  {}
func (noopLogger) Info(args ...interface{})           {}
func (noopLogger) Infof(t string, a ...interface{})   {}
func (noopLogger) Warn(args ...interface{})           {}
func (noopLogger) Warnf(t string, a ...interface{})   {}
func (noopLogger) Error(args ...interface{})          {}
func (noopLogger) Errorf(t string, a ...interface{})  {}
func (noopLogger) DPanic(args ...interface{})         {}
func (noopLogger) DPanicf(t string, a ...interface{}) {}
func (noopLogger) Fatal(args ...interface{})          {}
func (noopLogger) Fatalf(t string, a ...interface{})  {}
func (noopLogger) Sync() error                        { return nil }

type fakeVDB struct {
	mu       sync.Mutex
	docs     map[string]*vectordb.VectorDoc
	provider string

	indexed []string
	lastDoc *vectordb.VectorDoc
	kValues []int
}

func newFakeVDB() *fakeVDB {
	return &fakeVDB{docs: map[string]*vectordb.VectorDoc{}, provider: "pgvector"}
}

func (f *fakeVDB) Provider() string { return f.provider }
func (f *fakeVDB) Ping(context.Context) error {
	return nil
}
func (f *fakeVDB) EnsureIndex(context.Context, int) error { return nil }

func (f *fakeVDB) Index(_ context.Context, doc *vectordb.VectorDoc) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	id := fmt.Sprintf("id-%d", len(f.docs)+1)
	doc.ID = id
	f.docs[id] = doc
	f.indexed = append(f.indexed, doc.DocumentID)
	return id, nil
}

func (f *fakeVDB) Get(_ context.Context, id string) (*vectordb.VectorDoc, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	d, ok := f.docs[id]
	if !ok {
		return nil, vectordb.ErrNotFound
	}
	cp := *d
	return &cp, nil
}

func (f *fakeVDB) List(_ context.Context, _, _ int) ([]vectordb.VectorDoc, int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]vectordb.VectorDoc, 0, len(f.docs))
	for _, d := range f.docs {
		out = append(out, *d)
	}
	return out, int64(len(out)), nil
}

func (f *fakeVDB) Update(_ context.Context, id string, doc *vectordb.VectorDoc) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.lastDoc = doc
	if _, ok := f.docs[id]; !ok {
		return vectordb.ErrNotFound
	}
	f.docs[id] = doc
	return nil
}

func (f *fakeVDB) Delete(context.Context, string) error { return nil }

func (f *fakeVDB) SearchKNN(_ context.Context, _ []float32, k int) ([]vectordb.SearchHit, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.kValues = append(f.kValues, k)
	return nil, nil
}

type embedServer struct {
	calls int
	vec   []float32
}

func (e *embedServer) handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e.calls++
		w.Header().Set("Content-Type", "application/json")
		if e.vec == nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":{"message":"embed disabled in this test"}}`))
			return
		}
		out, _ := json.Marshal(map[string]interface{}{
			"data": []map[string]interface{}{{"embedding": e.vec}},
		})
		_, _ = w.Write(out)
	})
}

func newTestUseCase(t *testing.T, vdb vectordb.VectorDB, embed *embedServer) VectorDataUseCase {
	t.Helper()
	srv := httptest.NewServer(embed.handler())
	t.Cleanup(srv.Close)
	cfg := &config.Config{
		VectorDB: config.VectorDBConfig{Index: "vector_documents", Dims: 3},
		LLM:      config.LLMConfig{BaseURL: srv.URL, Model: "test", Dims: 3, Timeout: 5},
	}
	llmClient := llm.NewClient(&cfg.LLM, noopLogger{})
	return NewVectorDataUseCase(vdb, llmClient, cfg, noopLogger{})
}

func TestCreate_WithEmbedding_SkipsLLM(t *testing.T) {
	vdb := newFakeVDB()
	embed := &embedServer{vec: nil} // calls would fail
	u := newTestUseCase(t, vdb, embed)
	ctx := context.Background()

	resp, err := u.Create(ctx, &presenter.CreateDocumentRequest{
		DocumentID: "d1",
		Content:    "hello",
		Embedding:  []float32{1, 2, 3},
		SourceType: "iot",
		Metadata:   "m",
	})
	assert.NoError(t, err)
	assert.Equal(t, "d1", resp.DocumentID)
	assert.Equal(t, "id-1", resp.ID)
	assert.Equal(t, 3, resp.Dim)
	assert.Equal(t, 0, embed.calls, "embedding was provided; LLM must not be called")
	assert.Equal(t, []string{"d1"}, vdb.indexed)
}

func TestCreate_WithContentOnly_GeneratesEmbedding(t *testing.T) {
	vdb := newFakeVDB()
	embed := &embedServer{vec: []float32{0.5, 0.6, 0.7}}
	u := newTestUseCase(t, vdb, embed)

	resp, err := u.Create(context.Background(), &presenter.CreateDocumentRequest{
		DocumentID: "d2",
		Content:    "sensor reading",
		SourceType: "iot-manual",
	})
	assert.NoError(t, err)
	assert.Equal(t, "d2", resp.DocumentID)
	assert.Equal(t, 1, embed.calls)
	assert.Equal(t, 3, resp.Dim)
}

func TestUpdate_MetadataOnly_PreservesRest(t *testing.T) {
	vdb := newFakeVDB()
	embed := &embedServer{vec: nil} // must not be called
	u := newTestUseCase(t, vdb, embed)

	_, err := u.Create(context.Background(), &presenter.CreateDocumentRequest{
		DocumentID: "d1",
		Content:    "orig content",
		Embedding:  []float32{1, 2, 3},
		SourceType: "iot-manual",
		Metadata:   "old-m",
	})
	assert.NoError(t, err)

	resp, err := u.Update(context.Background(), "id-1", &presenter.UpdateDocumentRequest{
		Metadata: "new-m",
	})
	assert.NoError(t, err)
	assert.Equal(t, "new-m", resp.Metadata)
	assert.Equal(t, "orig content", resp.Content)
	assert.Equal(t, "iot-manual", resp.SourceType)
	assert.Equal(t, "d1", resp.DocumentID)

	vdb.mu.Lock()
	defer vdb.mu.Unlock()
	assert.Equal(t, "orig content", vdb.lastDoc.Content)
	assert.Equal(t, "new-m", vdb.lastDoc.Metadata)
	assert.Nil(t, vdb.lastDoc.Embedding, "metadata-only update must not regenerate embedding")
	assert.Equal(t, 0, embed.calls, "no LLM call expected")
}

func TestUpdate_ContentOnly_RegeneratesEmbedding(t *testing.T) {
	vdb := newFakeVDB()
	embed := &embedServer{vec: []float32{9, 9, 9}}
	u := newTestUseCase(t, vdb, embed)

	_, err := u.Create(context.Background(), &presenter.CreateDocumentRequest{
		DocumentID: "d1",
		Content:    "old",
		Embedding:  []float32{1, 2, 3},
		SourceType: "iot",
		Metadata:   "keep",
	})
	assert.NoError(t, err)

	resp, err := u.Update(context.Background(), "id-1", &presenter.UpdateDocumentRequest{
		Content: "new content",
	})
	assert.NoError(t, err)
	assert.Equal(t, "new content", resp.Content)
	assert.Equal(t, "keep", resp.Metadata, "metadata must be preserved")

	vdb.mu.Lock()
	defer vdb.mu.Unlock()
	assert.Equal(t, []float32{9, 9, 9}, vdb.lastDoc.Embedding)
	assert.Equal(t, "new content", vdb.lastDoc.Content)
	assert.Equal(t, 1, embed.calls)
}

func TestUpdate_NotFound_Propagates(t *testing.T) {
	vdb := newFakeVDB()
	embed := &embedServer{vec: []float32{1, 2, 3}}
	u := newTestUseCase(t, vdb, embed)

	_, err := u.Update(context.Background(), "missing", &presenter.UpdateDocumentRequest{
		Metadata: "x",
	})
	assert.ErrorIs(t, err, vectordb.ErrNotFound)
}

func TestSeed_NilLLM_ReturnsError(t *testing.T) {
	vdb := newFakeVDB()
	cfg := &config.Config{VectorDB: config.VectorDBConfig{Index: "vector_documents", Dims: 3}}
	u := NewVectorDataUseCase(vdb, nil, cfg, noopLogger{})

	_, err := u.Seed(context.Background(), &presenter.SeedRequest{Count: 1})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "llm client is nil")
}

func TestSearch_KIsCapped(t *testing.T) {
	vdb := newFakeVDB()
	embed := &embedServer{vec: []float32{0.1, 0.2, 0.3}}
	u := newTestUseCase(t, vdb, embed)

	t.Run("k above cap clamps to 100", func(t *testing.T) {
		_, err := u.Search(context.Background(), &presenter.SearchRequest{Query: "q", K: 999})
		assert.NoError(t, err)
		vdb.mu.Lock()
		assert.Equal(t, 100, vdb.kValues[len(vdb.kValues)-1])
		vdb.mu.Unlock()
	})

	t.Run("k default 10 when unset", func(t *testing.T) {
		_, err := u.Search(context.Background(), &presenter.SearchRequest{Query: "q"})
		assert.NoError(t, err)
		vdb.mu.Lock()
		assert.Equal(t, 10, vdb.kValues[len(vdb.kValues)-1])
		vdb.mu.Unlock()
	})
}

func TestSeed_Idempotent_SkipsExisting(t *testing.T) {
	vdb := newFakeVDB()
	embed := &embedServer{vec: []float32{0.1, 0.2, 0.3}}
	u := newTestUseCase(t, vdb, embed)

	ctx := context.Background()
	_, err := u.Create(ctx, &presenter.CreateDocumentRequest{
		DocumentID: "doc-temp-sensor",
		Content:    "existing",
		Embedding:  []float32{1, 2, 3},
	})
	assert.NoError(t, err)

	embed.calls = 0
	vdb.mu.Lock()
	indexedBefore := len(vdb.indexed)
	vdb.mu.Unlock()

	resp, err := u.Seed(ctx, &presenter.SeedRequest{Count: 5})
	assert.NoError(t, err)
	assert.Equal(t, 4, resp.Seeded, "existing doc-temp-sensor must be skipped")
	assert.Equal(t, 4, embed.calls)

	vdb.mu.Lock()
	defer vdb.mu.Unlock()
	for _, id := range vdb.indexed[indexedBefore:] {
		assert.NotEqual(t, "doc-temp-sensor", id, "duplicate must not be re-indexed")
	}
}
