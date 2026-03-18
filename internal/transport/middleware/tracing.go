package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/kongnakornna/golangapi/pkg/logger"
)

// TracingMiddleware request tracing middleware
func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get or generate request ID
		requestID := r.Header.Get("X-Request-ID")
		if requestID == " {
			requestID = middleware.GetReqID(r.Context())
			if requestID == " {
				requestID = fmt.Sprintf("%d", middleware.NextRequestID())
			}
		}

		// Get or generate trace ID
		traceID := r.Header.Get("X-Trace-ID")
		if traceID == " {
			traceID = r.Header.Get("X-B3-TraceId") // Support Zipkin B3 format
			if traceID == " {
				traceID = requestID // If no trace ID, use request ID
			}
		}

		// Get span ID (if exists)
		spanID := r.Header.Get("X-Span-ID")
		if spanID == " {
			spanID = r.Header.Get("X-B3-SpanId") // Support Zipkin B3 format
		}

		// Get parent span ID (if exists)
		parentSpanID := r.Header.Get("X-Parent-Span-ID")
		if parentSpanID == " {
			parentSpanID = r.Header.Get("X-B3-ParentSpanId") // Support Zipkin B3 format
		}

		// Set response headers
		w.Header().Set("X-Request-ID", requestID)
		w.Header().Set("X-Trace-ID", traceID)
		if spanID != " {
			w.Header().Set("X-Span-ID", spanID)
		}

		// Create context with tracing information
		ctx := r.Context()
		ctx = logger.WithRequestID(ctx, requestID)
		ctx = logger.WithTraceID(ctx, traceID)
		ctx = context.WithValue(ctx, "span_id", spanID)
		ctx = context.WithValue(ctx, "parent_span_id", parentSpanID)

		// Add tracing information to Chi context
		ctx = context.WithValue(ctx, middleware.RequestIDKey, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetTraceInfo gets tracing information from context
func GetTraceInfo(ctx context.Context) TraceInfo {
	return TraceInfo{
		RequestID:    logger.GetRequestID(ctx),
		TraceID:      logger.GetTraceID(ctx),
		SpanID:       getStringFromContext(ctx, "span_id"),
		ParentSpanID: getStringFromContext(ctx, "parent_span_id"),
	}
}

// TraceInfo tracing information
type TraceInfo struct {
	RequestID    string `json:"request_id"`
	TraceID      string `json:"trace_id"`
	SpanID       string `json:"span_id,omitempty"`
	ParentSpanID string `json:"parent_span_id,omitempty"`
}

// getStringFromContext gets a string value from context
func getStringFromContext(ctx context.Context, key string) string {
	if value := ctx.Value(key); value != nil {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return "
}

// RequestContextMiddleware request context middleware
func RequestContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create request context
		ctx := r.Context()

		// Add request method and path
		ctx = context.WithValue(ctx, "http_method", r.Method)
		ctx = context.WithValue(ctx, "http_path", r.URL.Path)
		ctx = context.WithValue(ctx, "http_query", r.URL.RawQuery)
		ctx = context.WithValue(ctx, "user_agent", r.UserAgent())
		ctx = context.WithValue(ctx, "remote_addr", r.RemoteAddr)

		// Add common request headers
		ctx = context.WithValue(ctx, "content_type", r.Header.Get("Content-Type"))
		ctx = context.WithValue(ctx, "accept", r.Header.Get("Accept"))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetHTTPRequestContext gets HTTP request context information
func GetHTTPRequestContext(ctx context.Context) HTTPRequestContext {
	return HTTPRequestContext{
		Method:      getStringFromContext(ctx, "http_method"),
		Path:        getStringFromContext(ctx, "http_path"),
		Query:       getStringFromContext(ctx, "http_query"),
		UserAgent:   getStringFromContext(ctx, "user_agent"),
		RemoteAddr:  getStringFromContext(ctx, "remote_addr"),
		ContentType: getStringFromContext(ctx, "content_type"),
		Accept:      getStringFromContext(ctx, "accept"),
	}
}

// HTTPRequestContext HTTP request context information
type HTTPRequestContext struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Query       string `json:"query,omitempty"`
	UserAgent   string `json:"user_agent"`
	RemoteAddr  string `json:"remote_addr"`
	ContentType string `json:"content_type,omitempty"`
	Accept      string `json:"accept,omitempty"`
}
