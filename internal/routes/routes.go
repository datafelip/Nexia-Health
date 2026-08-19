package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/felipelima/nexia-health/api/internal/handlers"
)

func Setup(pacienteHandler *handlers.PacienteHandler) http.Handler {
	router := chi.NewRouter()

	router.Get("/health", handlers.Health)

	router.Post("/pacientes", pacienteHandler.Create)

	router.Get("/pacientes", pacienteHandler.List)

	router.Get("/pacientes/{id}", pacienteHandler.GetByID)

	router.Put("/pacientes/{id}", pacienteHandler.Update)

	router.Delete("/pacientes/{id}", pacienteHandler.Delete)

	return router
}