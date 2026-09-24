package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"icmongolang/config"
	"icmongolang/internal/modules/settings"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/responses"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type settingsHandler struct {
	cfg    *config.Config
	uc     settings.SettingsUseCaseI
	logger logger.Logger
}

func CreateSettingsHandler(uc settings.SettingsUseCaseI, cfg *config.Config, logger logger.Logger) settings.Handlers {
	return &settingsHandler{cfg: cfg, uc: uc, logger: logger}
}

// ---- shared helpers ----

func listQueryFrom(r *http.Request) settings.ListQuery {
	q := r.URL.Query()
	vals := make(map[string]string, len(q))
	for k := range q {
		vals[k] = q.Get(k)
	}
	page, _ := strconv.Atoi(vals["page"])
	size, _ := strconv.Atoi(vals["pageSize"])
	return settings.NewListQuery(page, size, vals["sort"], vals)
}

func decodeBody(r *http.Request) (map[string]interface{}, error) {
	body := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	return body, nil
}

func filterKeys(body map[string]interface{}, allowed []string) map[string]interface{} {
	out := map[string]interface{}{}
	for _, k := range allowed {
		if v, ok := body[k]; ok && v != nil {
			out[k] = v
		}
	}
	return out
}

// coerceID turns numeric-looking ids into int so PG compares them against int columns.
func coerceID(s string) interface{} {
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return s
}

func ok(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.Respond(w, r, responses.CreateSuccessResponse(data))
}

func fail(w http.ResponseWriter, r *http.Request, err error) {
	render.Render(w, r, responses.CreateErrorResponse(err))
}

func badRequest(w http.ResponseWriter, r *http.Request, msg string) {
	render.Render(w, r, responses.CreateErrorResponse(
		responses.NewError(http.StatusBadRequest, msg)))
}

func response422(msg string) render.Renderer {
	return responses.CreateErrorResponse(responses.NewError(http.StatusUnprocessableEntity, msg))
}

func notFound(w http.ResponseWriter, r *http.Request, msg string) {
	render.Render(w, r, responses.CreateErrorResponse(
		responses.NewError(http.StatusNotFound, msg)))
}

func affected(w http.ResponseWriter, r *http.Request, n int64) {
	render.Respond(w, r, responses.CreateSuccessResponse(map[string]interface{}{"affected": n}))
}

func urlID(r *http.Request) string {
	return chi.URLParam(r, "id")
}

func nowPtr() time.Time {
	return time.Now()
}
