// internal/server/runtime.go
package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Runtime struct {
	Router   *chi.Mux
	shutdown []func(context.Context) error
}

func (r *Runtime) AddShutdown(fn func(context.Context) error) {
	if fn != nil {
		r.shutdown = append(r.shutdown, fn)
	}
}

func (r *Runtime) Shutdown(ctx context.Context) error {
	var firstErr error
	for i, fn := range r.shutdown {
		select {
		case <-ctx.Done():
			if firstErr == nil {
				firstErr = fmt.Errorf("shutdown timeout after %d hooks", i)
			}
			return firstErr
		default:
		}
		if err := fn(ctx); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (r *Runtime) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Router.ServeHTTP(w, req)
}
