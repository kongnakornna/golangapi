package template

import "net/http"

type Handlers interface {
	Create() http.HandlerFunc
	Get() http.HandlerFunc
	GetMulti() http.HandlerFunc
	Update() http.HandlerFunc
	Delete() http.HandlerFunc
}
