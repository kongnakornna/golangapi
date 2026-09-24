package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
)

type server struct {
	model     llms.Model
	chain     *chains.LLMChain
	modelName string
	log       *slog.Logger
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Messages []message `json:"messages"`
}

type chatResponse struct {
	Reply string `json:"reply"`
	Model string `json:"model"`
}

type generateRequest struct {
	Question string `json:"question"`
	Context  string `json:"context,omitempty"`
}

type generateResponse struct {
	Reply string `json:"reply"`
	Model string `json:"model"`
}

func (s *server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"model":  s.modelName,
	})
}

const systemPrompt = `You are an IoT operations assistant for the icmongolang platform.
Provide concise, factual answers about sensor data, alarms, MQTT topics, and platform events.
Use the provided context when relevant.
Be brief: keep answers under 150 words unless the user asks for detail.`

func (s *server) handleChat(w http.ResponseWriter, r *http.Request) {
	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	if len(req.Messages) == 0 {
		writeError(w, http.StatusBadRequest, "messages is required")
		return
	}

	content := make([]llms.MessageContent, 0, len(req.Messages)+1)
	if req.Messages[0].Role != "system" {
		content = append(content, llms.MessageContent{
			Role:  llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{llms.TextContent{Text: systemPrompt}},
		})
	}
	for _, m := range req.Messages {
		content = append(content, llms.MessageContent{
			Role:  llms.ChatMessageType(m.Role),
			Parts: []llms.ContentPart{llms.TextContent{Text: m.Content}},
		})
	}

	resp, err := s.model.GenerateContent(r.Context(), content,
		llms.WithMaxTokens(1024),
		llms.WithTemperature(0.2),
		llms.WithTopP(0.8),
	)
	if err != nil {
		s.log.Error("chat generation failed", "error", err)
		writeError(w, http.StatusInternalServerError, "generation failed")
		return
	}

	if len(resp.Choices) == 0 {
		writeError(w, http.StatusInternalServerError, "model returned no choices")
		return
	}

	writeJSON(w, http.StatusOK, chatResponse{
		Reply: resp.Choices[0].Content,
		Model: s.modelName,
	})
}

func (s *server) handleGenerate(w http.ResponseWriter, r *http.Request) {
	var req generateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	if req.Question == "" {
		writeError(w, http.StatusBadRequest, "question is required")
		return
	}

	out, err := chains.Call(r.Context(), s.chain, map[string]any{
		"context":  req.Context,
		"question": req.Question,
	}, chains.WithMaxTokens(1024))
	if err != nil {
		s.log.Error("chain call failed", "error", err)
		writeError(w, http.StatusInternalServerError, "generation failed")
		return
	}

	reply, _ := out["text"].(string)
	writeJSON(w, http.StatusOK, generateResponse{
		Reply: reply,
		Model: s.modelName,
	})
}

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(log)

	port := envOr("PORT", "8001")
	baseURL := envOr("OLLAMA_BASE_URL", "http://localhost:11434")
	modelName := envOr("OLLAMA_MODEL", "deepseek-r1:7b")
	timeout := envSeconds("AI_TIMEOUT", 180)

	model, err := ollama.New(
		ollama.WithModel(modelName),
		ollama.WithServerURL(baseURL),
		ollama.WithHTTPClient(&http.Client{Timeout: timeout}),
		ollama.WithRunnerNumCtx(2048),
	)
	if err != nil {
		log.Error("failed to init ollama model", "error", err)
		os.Exit(1)
	}

	prompt := prompts.NewPromptTemplate(
		systemPrompt+`

Context:
{{.context}}

Question:
{{.question}}`,
		[]string{"context", "question"},
	)

	chain := chains.NewLLMChain(model, prompt)
	s := &server{model: model, chain: chain, modelName: modelName, log: log}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("POST /v1/chat", s.handleChat)
	mux.HandleFunc("POST /v1/generate", s.handleGenerate)

	log.Info("ai service started", "port", port, "model", modelName, "ollama", baseURL)
	if err := http.ListenAndServe(":"+port, mux); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("http server error", "error", err)
		os.Exit(1)
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envSeconds(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
