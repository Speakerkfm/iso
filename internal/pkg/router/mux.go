package router

import (
	"github.com/go-chi/chi"
)

// ServeMux ...
type ServeMux chi.Router

// NewRouter ...
func NewRouter() *chi.Mux {
	return chi.NewRouter()
}
