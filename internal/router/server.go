package router

import (
	"time"

	"github.com/JesusEugenio/MitSys/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router           *chi.Mux
	inventoryHandler *handlers.InventoryHandler
}

func NewServer(invH *handlers.InventoryHandler) *Server {
	return &Server{
		router:           chi.NewRouter(),
		inventoryHandler: invH,
	}
}

func (s *Server) Setup() *chi.Mux {

	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	s.router.Use(middleware.Compress(5))
	s.router.Use(middleware.Timeout(60 * time.Second))

	s.router.Mount("/inventory", s.inventoryRoutes())

	return s.router
}
