package usecase

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"icmongolang/config"
	"icmongolang/internal/modules/elasticsearch/presenter"
	esclient "icmongolang/pkg/elasticsearch"
	"icmongolang/pkg/llm"

	"github.com/stretchr/testify/assert"
)

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

type esRecorder struct {
	docBodies  []string
	searchBody string
	lastIndex  string
}

func newFakeESServer(t *testing.T, rec *esRecorder) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/":
			_, _ = w.Write([]byte(`{"tagline":"You Know, for Search"}`))
		case r.Method == http.MethodPut && strings.Count(strings.Trim(r.URL.Path, "/"), "/") == 0:
			rec.lastIndex = strings.TrimPrefix(r.URL.Path, "/")
			_, _ = w.Write([]byte(`{"acknowledged":true}`))
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/_doc"):
			body, _ := io.ReadAll(r.Body)
			rec.docBodies = append(rec.docBodies, string(body))
			_, _ = w.Write([]byte(`{"_id":"doc-1","result":"created"}`))
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/_search"):
			body, _ := io.ReadAll(r.Body)
			rec.searchBody = string(body)
			_, _ = w.Write([]byte(`{"hits":{"total":{"value":1},"hits":[
				{"_id":"doc-1","_score":0.9,"_source":{"document_id":"a","content":"hello","source_type":"t","metadata":"meta-a"}}
			]}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"not found"}`))
		}
	})
	return httptest.NewServer(mux)
}

func newFakeLLMServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"embedding":[0.1,0.2,0.3]}]}`))
	}))
}

func newTestUC(t *testing.T, esSrv, llmSrv *httptest.Server) ElasticsearchUseCase {
	t.Helper()
	es, err := esclient.NewClient(&config.ElasticsearchConfig{
		Addresses: []string{esSrv.URL},
		Index:     "vector_embeddings",
		Timeout:   5,
	}, dummyLogger{})
	if err != nil {
		t.Fatal(err)
	}
	llmClient := llm.NewClient(&config.LLMConfig{BaseURL: llmSrv.URL, Model: "fake", Dims: 3, Timeout: 5}, dummyLogger{})
	return NewElasticsearchUseCase(es, llmClient, dummyLogger{})
}

func TestEmbedAndIndex_RoundTripsMetadataAndIndexName(t *testing.T) {
	rec := &esRecorder{}
	esSrv := newFakeESServer(t, rec)
	defer esSrv.Close()
	uc := newTestUC(t, esSrv, newFakeLLMServer(t))

	resp, err := uc.EmbedAndIndex(context.Background(), &presenter.EmbeddingRequest{
		DocumentID: "a",
		Content:    "hello",
		SourceType: "t",
		Metadata:   "meta-a",
	})
	assert.NoError(t, err)
	assert.True(t, resp.Indexed)
	assert.Equal(t, 3, resp.Dim)
	assert.Equal(t, "vector_embeddings", resp.Index)

	assert.Len(t, rec.docBodies, 1)
	var doc map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(rec.docBodies[0]), &doc))
	assert.Equal(t, "meta-a", doc["metadata"])
	assert.Equal(t, "a", doc["document_id"])
}

func TestIndexVector_RoundTripsMetadata(t *testing.T) {
	rec := &esRecorder{}
	esSrv := newFakeESServer(t, rec)
	defer esSrv.Close()
	uc := newTestUC(t, esSrv, newFakeLLMServer(t))

	resp, err := uc.IndexVector(context.Background(), &presenter.IndexVectorRequest{
		DocumentID: "a",
		Content:    "hello",
		Embedding:  []float32{1, 2, 3},
		SourceType: "t",
		Metadata:   "meta-a",
	})
	assert.NoError(t, err)
	assert.True(t, resp.Indexed)
	assert.Equal(t, "vector_embeddings", resp.Index)

	assert.Len(t, rec.docBodies, 1)
	var doc map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(rec.docBodies[0]), &doc))
	assert.Equal(t, "meta-a", doc["metadata"])
}

func TestSearchVector_MapsMetadataAndClampsK(t *testing.T) {
	rec := &esRecorder{}
	esSrv := newFakeESServer(t, rec)
	defer esSrv.Close()
	uc := newTestUC(t, esSrv, newFakeLLMServer(t))

	resp, err := uc.SearchVector(context.Background(), &presenter.SearchRequest{
		Query: "hello",
		K:     100000,
	})
	assert.NoError(t, err)
	assert.Equal(t, 100, resp.K)
	assert.Len(t, resp.Hits, 1)
	assert.Equal(t, "meta-a", resp.Hits[0].Metadata)
	assert.Equal(t, "a", resp.Hits[0].DocumentID)

	var q struct {
		Size int `json:"size"`
		KNN  struct {
			K int `json:"k"`
		} `json:"knn"`
	}
	assert.NoError(t, json.Unmarshal([]byte(rec.searchBody), &q))
	assert.Equal(t, 100, q.Size)
	assert.Equal(t, 100, q.KNN.K)
}

func TestCreateIndex_WithNilLLM_ReturnsError(t *testing.T) {
	esSrv := newFakeESServer(t, &esRecorder{})
	defer esSrv.Close()
	es, err := esclient.NewClient(&config.ElasticsearchConfig{
		Addresses: []string{esSrv.URL},
		Index:     "vector_embeddings",
		Timeout:   5,
	}, dummyLogger{})
	assert.NoError(t, err)
	uc := NewElasticsearchUseCase(es, nil, dummyLogger{})

	err = uc.CreateIndex(context.Background(), &presenter.CreateIndexRequest{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "llm")
}

func TestEmbedAndIndex_WithNilLLM_ReturnsError(t *testing.T) {
	esSrv := newFakeESServer(t, &esRecorder{})
	defer esSrv.Close()
	es, err := esclient.NewClient(&config.ElasticsearchConfig{
		Addresses: []string{esSrv.URL},
		Index:     "vector_embeddings",
		Timeout:   5,
	}, dummyLogger{})
	assert.NoError(t, err)
	uc := NewElasticsearchUseCase(es, nil, dummyLogger{})

	_, err = uc.EmbedAndIndex(context.Background(), &presenter.EmbeddingRequest{Content: "x"})
	assert.Error(t, err)
}
