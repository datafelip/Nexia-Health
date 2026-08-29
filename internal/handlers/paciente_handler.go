package handlers

import (
	"encoding/json"
	"net/http"
	"errors"
	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/felipelima/nexia-health/api/internal/repositories"
	"github.com/felipelima/nexia-health/api/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
		ErrorResponseJSON(w, http.StatusBadRequest, "ID do paciente inválido")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrPacienteNotFound) {
			ErrorResponseJSON(w, http.StatusNotFound, "paciente não encontrado")
			return
		}

		ErrorResponseJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
