package router

import (
	"github.com/go-chi/chi/v5"
)

func (s *Server) inventoryRoutes() chi.Router {
	r := chi.NewRouter()

	r.Route("/new", func(r chi.Router) {
		r.Post("/", s.inventoryHandler.CreateProduct)
		r.Get("/", s.inventoryHandler.CreateProductForm)
	})

	return r
}
