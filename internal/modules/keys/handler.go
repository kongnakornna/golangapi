package keys

import "net/http"

type Handlers interface {
	CreateAPIKey() http.HandlerFunc
	ListAPIKeys() http.HandlerFunc
	GetAPIKey() http.HandlerFunc
	RotateAPIKey() http.HandlerFunc
	RevokeAPIKey() http.HandlerFunc
}