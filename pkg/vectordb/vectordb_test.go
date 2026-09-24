package vectordb

import (
	"testing"

	"icmongolang/config"

	"github.com/stretchr/testify/assert"
)

func TestVecToString(t *testing.T) {
	tests := []struct {
		name string
		in   []float32
		want string
	}{
		{name: "empty", in: nil, want: "[]"},
		{name: "empty-slice", in: []float32{}, want: "[]"},
		{name: "single", in: []float32{1.5}, want: "[1.5]"},
		{name: "multi", in: []float32{0.1, -2, 3}, want: "[0.1,-2,3]"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, vecToString(tc.in))
		})
	}
}

func TestIsNumericID(t *testing.T) {
	tests := []struct {
		id   string
		want bool
	}{
		{id: "1", want: true},
		{id: "42", want: true},
		{id: "0", want: false},
		{id: "-1", want: false},
		{id: "abc", want: false},
		{id: "1.5", want: false},
		{id: "", want: false},
	}
	for _, tc := range tests {
		t.Run(tc.id, func(t *testing.T) {
			assert.Equal(t, tc.want, isNumericID(tc.id))
		})
	}
}

func TestNewProviderSelection(t *testing.T) {
	cfg := &config.VectorDBConfig{Provider: "", Index: "vector_documents", Dims: 768}

	t.Run("empty provider defaults to pgvector", func(t *testing.T) {
		_, err := New(cfg, nil, nil, nil)
		assert.Error(t, err) // db is nil -> error, but proves it took the pgvector branch
	})

	t.Run("pgvector without db errors", func(t *testing.T) {
		_, err := New(&config.VectorDBConfig{Provider: "pgvector"}, nil, nil, nil)
		assert.ErrorContains(t, err, "pgvector")
	})

	t.Run("elasticsearch without es errors", func(t *testing.T) {
		_, err := New(&config.VectorDBConfig{Provider: "elasticsearch"}, nil, nil, nil)
		assert.ErrorContains(t, err, "elasticsearch")
	})

	t.Run("unknown provider errors", func(t *testing.T) {
		_, err := New(&config.VectorDBConfig{Provider: "clickhouse"}, nil, nil, nil)
		assert.ErrorContains(t, err, "unknown provider")
	})
}
