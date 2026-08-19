package handlers

import (
	"encoding/json"
	"github.com/felipelima/nexia-health/api/internal/models"
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(paciente); err != nil {
		http.Error(w, "erro ao gerar resposta", http.StatusInternalServerError)
	}
}

func (h *PacienteHandler) List(w http.ResponseWriter, r *http.Request) {
	pacientes, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, "erro ao buscar pacientes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(pacientes); err != nil {
		http.Error(w, "erro ao gerar resposta", http.StatusInternalServerError)
	}
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(paciente); err != nil {
		http.Error(w, "erro ao gerar resposta", http.StatusInternalServerError)
	}
}

func (h *PacienteHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "ID do paciente inválido", http.StatusBadRequest)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(pacienteAtualizado); err != nil {
		http.Error(w, "erro ao gerar resposta", http.StatusInternalServerError)
	}
}

func (h *PacienteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "ID do paciente inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if err.Error() == "paciente não encontrado" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}