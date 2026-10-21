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

type ProfissionalHandler struct {
	service *services.ProfissionalService
}

func NewProfissionalHandler(
	service *services.ProfissionalService,
) *ProfissionalHandler {
	return &ProfissionalHandler{service: service}
}

func (h *ProfissionalHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var profissional models.Profissional

	if err := json.NewDecoder(r.Body).Decode(&profissional); err != nil {
		writeProfissionalJSON(w, http.StatusBadRequest, map[string]string{
			"erro": "JSON inválido",
		})
		return
	}

	if err := h.service.Create(r.Context(), &profissional); err != nil {
		writeProfissionalError(w, err)
		return
	}

	writeProfissionalJSON(w, http.StatusCreated, profissional)
}

func (h *ProfissionalHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {
	cpf := r.URL.Query().Get("cpf")

	profissionais, err := h.service.List(r.Context(), cpf)
	if err != nil {
		writeProfissionalError(w, err)
		return
	}

	writeProfissionalJSON(w, http.StatusOK, profissionais)
}

func (h *ProfissionalHandler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeProfissionalJSON(w, http.StatusBadRequest, map[string]string{
			"erro": "ID do profissional inválido",
		})
		return
	}

	profissional, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeProfissionalError(w, err)
		return
	}

	writeProfissionalJSON(w, http.StatusOK, profissional)
}

func (h *ProfissionalHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeProfissionalJSON(w, http.StatusBadRequest, map[string]string{
			"erro": "ID do profissional inválido",
		})
		return
	}

	var profissional models.Profissional

	if err := json.NewDecoder(r.Body).Decode(&profissional); err != nil {
		writeProfissionalJSON(w, http.StatusBadRequest, map[string]string{
			"erro": "JSON inválido",
		})
		return
	}

	profissional.ID = id

	if err := h.service.Update(r.Context(), &profissional); err != nil {
		writeProfissionalError(w, err)
		return
	}

	writeProfissionalJSON(w, http.StatusOK, profissional)
}

func (h *ProfissionalHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeProfissionalJSON(w, http.StatusBadRequest, map[string]string{
			"erro": "ID do profissional inválido",
		})
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeProfissionalError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeProfissionalError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrProfissionalDadosInvalidos):
		writeProfissionalJSON(w, http.StatusBadRequest, map[string]string{
			"erro": err.Error(),
		})

	case errors.Is(err, services.ErrProfissionalIDInvalido):
		writeProfissionalJSON(w, http.StatusBadRequest, map[string]string{
			"erro": err.Error(),
		})

	case errors.Is(err, repositories.ErrProfissionalNaoEncontrado):
		writeProfissionalJSON(w, http.StatusNotFound, map[string]string{
			"erro": err.Error(),
		})

	case errors.Is(err, repositories.ErrProfissionalDuplicado):
		writeProfissionalJSON(w, http.StatusConflict, map[string]string{
			"erro": err.Error(),
		})

	case errors.Is(err, repositories.ErrProfissionalReferencia):
		writeProfissionalJSON(w, http.StatusBadRequest, map[string]string{
			"erro": err.Error(),
		})

	default:
		writeProfissionalJSON(w, http.StatusInternalServerError, map[string]string{
			"erro": "Erro interno ao processar a solicitação",
		})
	}
}

func writeProfissionalJSON(
	w http.ResponseWriter,
	status int,
	data any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
