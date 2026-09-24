package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"icmongolang/config"
	"icmongolang/pkg/logger"
)

// Client is an OpenAI-compatible embedding client.
// Works with OpenAI, Azure OpenAI, Ollama, local servers, etc.
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	model      string
	dims       int
	logger     logger.Logger
}

// NewClient creates an LLM embedding client.
func NewClient(cfg *config.LLMConfig, log logger.Logger) *Client {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 30
	}

	dims := cfg.Dims
	if dims <= 0 {
		dims = 1536
	}

	return &Client{
		httpClient: &http.Client{Timeout: time.Duration(timeout) * time.Second},
		baseURL:    strings.TrimSuffix(cfg.BaseURL, "/"),
		apiKey:     cfg.APIKey,
		model:      cfg.Model,
		dims:       dims,
		logger:     log,
	}
}

// Dims returns the configured embedding dimension.
func (c *Client) Dims() int {
	return c.dims
}

// Model returns the configured embedding model name.
func (c *Client) Model() string {
	return c.model
}

// Embed converts a text into a vector embedding using the configured model.
func (c *Client) Embed(ctx context.Context, text string) ([]float32, error) {
	body, err := json.Marshal(map[string]interface{}{
		"model": c.model,
		"input": text,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("llm request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("llm embeddings returned %d: %s", resp.StatusCode, string(respBody))
	}

	var out struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("llm response decode: %w", err)
	}

	if out.Error != nil {
		return nil, fmt.Errorf("llm error: %s", out.Error.Message)
	}

	if len(out.Data) == 0 || len(out.Data[0].Embedding) == 0 {
		return nil, fmt.Errorf("llm returned empty embedding")
	}

	return out.Data[0].Embedding, nil
}
