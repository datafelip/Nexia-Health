package routes

import (
	"net/http"

	"github.com/felipelima/nexia-health/api/internal/handlers"
	"github.com/go-chi/chi/v5"
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
	profissionalHandler *handlers.ProfissionalHandler,
	exameHandler *handlers.ExameHandler,
	solicitacaoExameHandler *handlers.SolicitacaoExameHandler,
	aiHandler *handlers.AIHandler,
) http.Handler {

	r := chi.NewRouter()

	// Health Check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
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

		// Reclassificação de risco.
		r.Patch("/{id}/reclassificar", triagemHandler.Reclassificar)

		// Cancelamento de triagem com preservação do histórico.
		// A operação também atualiza e reorganiza a fila.
		r.Patch("/{id}/cancelar", triagemHandler.Cancelar)

		// Edição direta não permitida.
		r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
			handlers.ErrorResponseJSON(
				w,
				http.StatusMethodNotAllowed,
				"A edição direta de triagens não é permitida",
			)
		})

		// Exclusão física não permitida.
		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
			handlers.ErrorResponseJSON(
				w,
				http.StatusMethodNotAllowed,
				"A exclusão de triagens não é permitida",
			)
		})
	})

	// Filas
	r.Route("/filas", func(r chi.Router) {
		r.Get("/", filaHandler.List)
		r.Get("/{id}", filaHandler.GetByID)

		r.Post("/chamar-proximo", filaHandler.ChamarProximo)
		r.Patch("/{id}/cancelar", filaHandler.Cancelar)

		// Criação direta de fila não permitida.
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			handlers.ErrorResponseJSON(
				w,
				http.StatusMethodNotAllowed,
				"A entrada na fila deve ocorrer pela triagem",
			)
		})

		// Edição direta de fila não permitida.
		r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
			handlers.ErrorResponseJSON(
				w,
				http.StatusMethodNotAllowed,
				"A fila é gerenciada pelas operações de atendimento",
			)
		})

		// Exclusão direta de fila não permitida.
		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
			handlers.ErrorResponseJSON(
				w,
				http.StatusMethodNotAllowed,
				"Exclusão direta da fila não permitida",
			)
		})
	})

	// Atendimentos
	r.Route("/atendimentos", func(r chi.Router) {
		r.Post("/", atendimentoHandler.Create)
		r.Get("/", atendimentoHandler.List)
		r.Get("/{id}", atendimentoHandler.GetByID)
		r.Put("/{id}", atendimentoHandler.Update)

		// Exclusão de atendimentos não permitida.
		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
			handlers.ErrorResponseJSON(
				w,
				http.StatusMethodNotAllowed,
				"A exclusão de atendimentos não é permitida",
			)
		})
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

	// Profissionais
	r.Route("/profissionais", func(r chi.Router) {
		r.Post("/", profissionalHandler.Create)
		r.Get("/", profissionalHandler.List)
		r.Get("/{id}", profissionalHandler.GetByID)
		r.Put("/{id}", profissionalHandler.Update)
		r.Delete("/{id}", profissionalHandler.Delete)
	})

	// Exames
	r.Route("/exames", func(r chi.Router) {
		r.Post("/", exameHandler.Create)
		r.Get("/", exameHandler.List)
		r.Get("/{id}", exameHandler.GetByID)
		r.Put("/{id}", exameHandler.Update)
		r.Delete("/{id}", exameHandler.Delete)
	})

	// Solicitações de Exames
	r.Route("/solicitacoes-exame", func(r chi.Router) {
		r.Post("/", solicitacaoExameHandler.Create)
		r.Get("/", solicitacaoExameHandler.List)
		r.Get("/{id}", solicitacaoExameHandler.GetByID)
		r.Put("/{id}", solicitacaoExameHandler.Update)
		r.Delete("/{id}", solicitacaoExameHandler.Delete)
	})

	// Nexia AI - Inteligência Artificial
	r.Route("/ai", func(r chi.Router) {
		r.Post("/prever-demanda", aiHandler.PreverDemanda)
	})

	return r
}
