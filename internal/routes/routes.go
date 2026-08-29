package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/felipelima/nexia-health/api/internal/handlers"
)

func Setup(
	pacienteHandler *handlers.PacienteHandler,
	triagemHandler *handlers.TriagemHandler,
) http.Handler {

	router := chi.NewRouter()

	// Health Check
	router.Get("/health", handlers.Health)

	// =========================
	// PACIENTES
	// =========================

	router.Post("/pacientes", pacienteHandler.Create)
	router.Get("/pacientes", pacienteHandler.List)
	router.Get("/pacientes/{id}", pacienteHandler.GetByID)
	router.Put("/pacientes/{id}", pacienteHandler.Update)
	router.Delete("/pacientes/{id}", pacienteHandler.Delete)

	// =========================
	// TRIAGENS
	// =========================

	router.Post("/triagens", triagemHandler.Create)
	router.Get("/triagens", triagemHandler.List)
	router.Get("/triagens/{id}", triagemHandler.GetByID)
	router.Put("/triagens/{id}", triagemHandler.Update)
	router.Delete("/triagens/{id}", triagemHandler.Delete)

	return router
}
