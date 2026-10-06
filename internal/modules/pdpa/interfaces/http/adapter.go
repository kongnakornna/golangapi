package http

import (
	"net/http"

	"github.com/google/uuid"

	"icmongolang/internal/middleware"
)

func currentUserID(r *http.Request) (uuid.UUID, bool) {
	u, err := middleware.GetUserFromCtx(r.Context())
	if err != nil || u == nil {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(u.ID.String())
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}
