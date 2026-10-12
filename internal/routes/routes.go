package routes

import (
	"github.com/felipelima/nexia-health/api/internal/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Setup(
	pacienteHandler *handlers.PacienteHandler,
	triagemHandler *handlers.TriagemHandler,
	filaHandler *handlers.FilaHandler,
	atendimentoHandler *handlers.AtendimentoHandler,
	agendamentoHandler *handlers.AgendamentoHandler,
	especialidadeHandler *handlers.EspecialidadeHandler,
	setorHandler *handlers.SetorHandler,
	usuarioHandler *handlers.UsuarioHandler,
) http.Handler {

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Pacientes
	r.Route("/pacientes", func(r chi.Router) {
		r.Post("/", pacienteHandler.Create)
		r.Get("/", pacienteHandler.List)
		r.Get("/{id}", pacienteHandler.GetByID)
		r.Put("/{id}", pacienteHandler.Update)
		r.Delete("/{id}", pacienteHandler.Delete)
	})

	// Triagens
	r.Route("/triagens", func(r chi.Router) {
		r.Post("/", triagemHandler.Create)
		r.Get("/", triagemHandler.List)
		r.Get("/{id}", triagemHandler.GetByID)
		r.Put("/{id}", triagemHandler.Update)
		r.Delete("/{id}", triagemHandler.Delete)
	})

	// Filas
	r.Route("/filas", func(r chi.Router) {
		r.Post("/", filaHandler.Create)
		r.Get("/", filaHandler.List)
		r.Get("/{id}", filaHandler.GetByID)
		r.Put("/{id}", filaHandler.Update)
		r.Delete("/{id}", filaHandler.Delete)
	})

	// Atendimentos
	r.Route("/atendimentos", func(r chi.Router) {
		r.Post("/", atendimentoHandler.Create)
		r.Get("/", atendimentoHandler.List)
		r.Get("/{id}", atendimentoHandler.GetByID)
		r.Put("/{id}", atendimentoHandler.Update)
		r.Delete("/{id}", atendimentoHandler.Delete)
	})

	// Agendamentos
	r.Route("/agendamentos", func(r chi.Router) {
		r.Post("/", agendamentoHandler.Create)
		r.Get("/", agendamentoHandler.List)
		r.Get("/{id}", agendamentoHandler.GetByID)
		r.Put("/{id}", agendamentoHandler.Update)
		r.Delete("/{id}", agendamentoHandler.Delete)
	})

	// Especialidades
	r.Route("/especialidades", func(r chi.Router) {
		r.Post("/", especialidadeHandler.Create)
		r.Get("/", especialidadeHandler.List)
		r.Get("/{id}", especialidadeHandler.GetByID)
		r.Put("/{id}", especialidadeHandler.Update)
		r.Delete("/{id}", especialidadeHandler.Delete)
	})

	// Setores
	r.Route("/setores", func(r chi.Router) {
		r.Post("/", setorHandler.Create)
		r.Get("/", setorHandler.List)
		r.Get("/{id}", setorHandler.GetByID)
		r.Put("/{id}", setorHandler.Update)
		r.Delete("/{id}", setorHandler.Delete)
	})

	// Usuários
	r.Route("/usuarios", func(r chi.Router) {
		r.Post("/", usuarioHandler.Create)
		r.Get("/", usuarioHandler.List)
		r.Get("/{id}", usuarioHandler.GetByID)
		r.Put("/{id}", usuarioHandler.Update)
		r.Delete("/{id}", usuarioHandler.Delete)
	})
	return r
}
