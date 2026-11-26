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

	r.Route("/update-price", func(r chi.Router) {
		r.Get("/", s.inventoryHandler.UpdateProductPriceForm)
		r.Patch("/", s.inventoryHandler.UpdateProductPrice)
	})

	r.Route("/catalog", func(r chi.Router) {
		r.Get("/", s.inventoryHandler.GetCatalog)
	})

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", s.inventoryHandler.ViewProduct)
	})

	r.Route("/update/{id}", func(r chi.Router) {
		r.Get("/", s.inventoryHandler.UpdateProductForm)
		r.Put("/", s.inventoryHandler.UpdateProduct)
	})

	return r
}
