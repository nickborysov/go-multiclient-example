package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nickborysov/go-multiclient-example/internal/dependency"
)

type Router struct {
	service dependency.Service
	*mux.Router
}

func NewRouter(service dependency.Service) *Router {
	r := &Router{
		Router:  mux.NewRouter(),
		service: service,
	}

	r.RegisterRoutes()

	return r
}

func (r *Router) RegisterRoutes() {
	r.Path("/example").
		Methods(http.MethodGet).
		HandlerFunc(r.HandleExample)
}

func (r *Router) HandleExample(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp, err := r.service.GetTestResponse()
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Fatal(err)
	}
}
