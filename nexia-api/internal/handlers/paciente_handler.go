package handlers

import (
	"encoding/json"
	"errors"
	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/felipelima/nexia-health/api/internal/repositories"
	"github.com/felipelima/nexia-health/api/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

type PacienteHandler struct {
	service *services.PacienteService
}

func NewPacienteHandler(service *services.PacienteService) *PacienteHandler {
	return &PacienteHandler{
		service: service,
	}
}

func (h *PacienteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var paciente models.Paciente

	if err := json.NewDecoder(r.Body).Decode(&paciente); err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.service.Create(r.Context(), &paciente); err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	JSONResponse(w, http.StatusCreated, paciente)
}

func (h *PacienteHandler) List(w http.ResponseWriter, r *http.Request) {
	pacientes, err := h.service.List(r.Context())
	if err != nil {
		ErrorResponseJSON(w, http.StatusInternalServerError, "erro ao buscar pacientes")
		return
	}

	JSONResponse(w, http.StatusOK, pacientes)
}

func (h *PacienteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "ID do paciente inválido")
		return
	}

	paciente, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		ErrorResponseJSON(w, http.StatusNotFound, "paciente não encontrado")
		return
	}

	JSONResponse(w, http.StatusOK, paciente)
}

func (h *PacienteHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "ID do paciente inválido")
		return
	}

	var paciente models.Paciente

	if err := json.NewDecoder(r.Body).Decode(&paciente); err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	paciente.ID = id

	pacienteAtualizado, err := h.service.Update(r.Context(), &paciente)
	if err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, pacienteAtualizado)
}

func (h *PacienteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponseJSON(
			w,
			http.StatusBadRequest,
			"ID do paciente inválido",
		)
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrPacienteNotFound):
			ErrorResponseJSON(
				w,
				http.StatusNotFound,
				"Paciente não encontrado ou já excluído",
			)

		case errors.Is(err, repositories.ErrPacienteFilaAtiva):
			ErrorResponseJSON(
				w,
				http.StatusConflict,
				"Não é possível excluir este paciente porque ele está em uma fila ativa. Finalize ou cancele sua participação na fila antes de continuar.",
			)

		case errors.Is(err, repositories.ErrPacienteAtendimentoAtivo):
			ErrorResponseJSON(
				w,
				http.StatusConflict,
				"Não é possível excluir este paciente porque ele possui um atendimento em andamento. Finalize o atendimento antes de continuar.",
			)

		default:
			ErrorResponseJSON(
				w,
				http.StatusInternalServerError,
				"Erro interno ao excluir paciente",
			)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
