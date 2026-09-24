package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ollamaUp = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ollama_up",
		Help: "Whether the Ollama API is reachable (1) or not (0)",
	})
	ollamaScrapeDuration = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ollama_scrape_duration_seconds",
		Help: "Duration of the last Ollama API scrape in seconds",
	})
	ollamaVersionInfo = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ollama_version_info",
		Help: "Ollama server version (1 if present)",
	}, []string{"version"})
	ollamaLoadedModels = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ollama_loaded_models",
		Help: "Models currently loaded in Ollama (1 if loaded)",
	}, []string{"model", "family", "parameter_size", "quantization"})
	ollamaModelSizeBytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ollama_model_size_bytes",
		Help: "Size of the model on disk in bytes",
	}, []string{"model"})
	ollamaModelVRAMBytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ollama_model_vram_bytes",
		Help: "VRAM used by the loaded model in bytes",
	}, []string{"model"})
	ollamaCatalogModels = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ollama_catalog_models",
		Help: "Models available in the Ollama catalog (1 if present)",
	}, []string{"model"})
	ollamaModelModifiedSeconds = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ollama_model_modified_seconds",
		Help: "Last modified time of the model as a unix timestamp",
	}, []string{"model"})
)

type psResponse struct {
	Models []struct {
		Name      string `json:"name"`
		Size      int64  `json:"size"`
		SizeVRAM  int64  `json:"size_vram"`
		ExpiresAt string `json:"expires_at"`
		Details   struct {
			Family            string `json:"family"`
			ParameterSize     string `json:"parameter_size"`
			QuantizationLevel string `json:"quantization_level"`
		} `json:"details"`
	} `json:"models"`
}

type tagsResponse struct {
	Models []struct {
		Name       string `json:"name"`
		Size       int64  `json:"size"`
		ModifiedAt string `json:"modified_at"`
	} `json:"models"`
}

type versionResponse struct {
	Version string `json:"version"`
}

func main() {
	listen := envOrDefault("LISTEN_ADDR", ":9309")
	baseURL := strings.TrimRight(envOrDefault("OLLAMA_URL", "http://localhost:11434"), "/")
	interval := 15 * time.Second

	flag.StringVar(&listen, "listen", listen, "address to expose /metrics on")
	flag.StringVar(&baseURL, "ollama.url", baseURL, "Ollama base URL")
	flag.DurationVar(&interval, "interval", interval, "scrape interval")
	flag.Parse()

	go func() {
		for {
			scrape(baseURL)
			time.Sleep(interval)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("ollama-exporter listening on %s, scraping %s every %s", listen, baseURL, interval)
	log.Fatal(http.ListenAndServe(listen, nil))
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func scrape(baseURL string) {
	start := time.Now()
	client := &http.Client{Timeout: 5 * time.Second}

	resetMetrics()

	up := 1.0
	defer func() {
		ollamaUp.Set(up)
		ollamaScrapeDuration.Set(time.Since(start).Seconds())
	}()

	var ver versionResponse
	if err := getJSON(client, baseURL+"/api/version", &ver); err != nil {
		log.Printf("ollama version scrape failed: %v", err)
		up = 0
		return
	}
	if ver.Version != "" {
		ollamaVersionInfo.WithLabelValues(ver.Version).Set(1)
	}

	var ps psResponse
	if err := getJSON(client, baseURL+"/api/ps", &ps); err != nil {
		log.Printf("ollama ps scrape failed: %v", err)
		up = 0
		return
	}
	for _, m := range ps.Models {
		ollamaLoadedModels.WithLabelValues(m.Name, m.Details.Family, m.Details.ParameterSize, m.Details.QuantizationLevel).Set(1)
		ollamaModelSizeBytes.WithLabelValues(m.Name).Set(float64(m.Size))
		ollamaModelVRAMBytes.WithLabelValues(m.Name).Set(float64(m.SizeVRAM))
	}

	var tags tagsResponse
	if err := getJSON(client, baseURL+"/api/tags", &tags); err != nil {
		log.Printf("ollama tags scrape failed: %v", err)
		up = 0
		return
	}
	for _, m := range tags.Models {
		ollamaCatalogModels.WithLabelValues(m.Name).Set(1)
		if t, err := time.Parse(time.RFC3339Nano, m.ModifiedAt); err == nil {
			ollamaModelModifiedSeconds.WithLabelValues(m.Name).Set(float64(t.Unix()))
		}
	}
}

func resetMetrics() {
	ollamaVersionInfo.Reset()
	ollamaLoadedModels.Reset()
	ollamaModelSizeBytes.Reset()
	ollamaModelVRAMBytes.Reset()
	ollamaCatalogModels.Reset()
	ollamaModelModifiedSeconds.Reset()
}

func getJSON(client *http.Client, url string, v interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %s", resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(v)
}
