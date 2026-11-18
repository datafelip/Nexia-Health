package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/felipelima/nexia-health/api/internal/repositories"
	"github.com/felipelima/nexia-health/api/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ExameHandler struct {
	service *services.ExameService
}

func NewExameHandler(
	service *services.ExameService,
) *ExameHandler {
	return &ExameHandler{service: service}
}

func (h *ExameHandler) Create(w http.ResponseWriter, r *http.Request) {
	var exame models.Exame

	if err := json.NewDecoder(r.Body).Decode(&exame); err != nil {
		writeExameJSON(w, http.StatusBadRequest, map[string]string{
			"erro": "JSON inválido",
		})
		return
	}

	if err := h.service.Create(r.Context(), &exame); err != nil {
		writeExameError(w, err)
		return
	}

	writeExameJSON(w, http.StatusCreated, exame)
}

func (h *ExameHandler) List(w http.ResponseWriter, r *http.Request) {
	nome := r.URL.Query().Get("nome")

	exames, err := h.service.List(r.Context(), nome)
	if err != nil {
		writeExameError(w, err)
		return
	}

	writeExameJSON(w, http.StatusOK, exames)
}

func (h *ExameHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeExameJSON(w, http.StatusBadRequest, map[string]string{
			"erro": "ID do exame inválido",
		})
		return
	}

	exame, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeExameError(w, err)
		return
	}

	writeExameJSON(w, http.StatusOK, exame)
}

func (h *ExameHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeExameJSON(w, http.StatusBadRequest, map[string]string{
			"erro": "ID do exame inválido",
		})
		return
	}

	var exame models.Exame

	if err := json.NewDecoder(r.Body).Decode(&exame); err != nil {
		writeExameJSON(w, http.StatusBadRequest, map[string]string{
			"erro": "JSON inválido",
		})
		return
	}

	exame.ID = id

	if err := h.service.Update(r.Context(), &exame); err != nil {
		writeExameError(w, err)
		return
	}

	writeExameJSON(w, http.StatusOK, exame)
}

func (h *ExameHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeExameJSON(w, http.StatusBadRequest, map[string]string{
			"erro": "ID do exame inválido",
		})
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeExameError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeExameError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrExameDadosInvalidos),
		errors.Is(err, services.ErrExameIDInvalido):
		writeExameJSON(w, http.StatusBadRequest, map[string]string{
			"erro": err.Error(),
		})

	case errors.Is(err, repositories.ErrExameNaoEncontrado):
		writeExameJSON(w, http.StatusNotFound, map[string]string{
			"erro": err.Error(),
		})

	case errors.Is(err, repositories.ErrExameDuplicado):
		writeExameJSON(w, http.StatusConflict, map[string]string{
			"erro": err.Error(),
		})

	case errors.Is(err, repositories.ErrExameEmUso):
		writeExameJSON(w, http.StatusConflict, map[string]string{
			"erro": err.Error(),
		})

	default:
		writeExameJSON(w, http.StatusInternalServerError, map[string]string{
			"erro": "Erro interno ao processar a solicitação",
		})
	}
}

func writeExameJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
