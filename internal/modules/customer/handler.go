package customer

import "net/http"

type Handlers interface {
	Create() func(w http.ResponseWriter, r *http.Request)
	Get() func(w http.ResponseWriter, r *http.Request)
	GetMulti() func(w http.ResponseWriter, r *http.Request)
	Delete() func(w http.ResponseWriter, r *http.Request)
	Update() func(w http.ResponseWriter, r *http.Request)
	CreateCar() func(w http.ResponseWriter, r *http.Request)
	GetCar() func(w http.ResponseWriter, r *http.Request)
	ListCars() func(w http.ResponseWriter, r *http.Request)
	UpdateCar() func(w http.ResponseWriter, r *http.Request)
	DeleteCar() func(w http.ResponseWriter, r *http.Request)
}
