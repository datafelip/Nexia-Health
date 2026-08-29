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

type TriagemHandler struct {
	service *services.TriagemService
}

func NewTriagemHandler(service *services.TriagemService) *TriagemHandler {
	return &TriagemHandler{
		service: service,
	}
}

func (h *TriagemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var triagem models.Triagem

	if err := json.NewDecoder(r.Body).Decode(&triagem); err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.service.Create(r.Context(), &triagem); err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	JSONResponse(w, http.StatusCreated, triagem)
}

func (h *TriagemHandler) List(w http.ResponseWriter, r *http.Request) {
	triagens, err := h.service.List(r.Context())
	if err != nil {
		ErrorResponseJSON(w, http.StatusInternalServerError, "erro ao buscar triagens")
		return
	}

	JSONResponse(w, http.StatusOK, triagens)
}

func (h *TriagemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "ID da triagem inválido")
		return
	}

	triagem, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrTriagemNotFound) {
			ErrorResponseJSON(w, http.StatusNotFound, "triagem não encontrada")
			return
		}

		ErrorResponseJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, triagem)
}

func (h *TriagemHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "ID da triagem inválido")
		return
	}

	var triagem models.Triagem

	if err := json.NewDecoder(r.Body).Decode(&triagem); err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	triagem.ID = id

	triagemAtualizada, err := h.service.Update(r.Context(), &triagem)
	if err != nil {
		if errors.Is(err, repositories.ErrTriagemNotFound) {
			ErrorResponseJSON(w, http.StatusNotFound, "triagem não encontrada")
			return
		}

		ErrorResponseJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, triagemAtualizada)
}

func (h *TriagemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "ID da triagem inválido")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrTriagemNotFound) {
			ErrorResponseJSON(w, http.StatusNotFound, "triagem não encontrada")
			return
		}

		ErrorResponseJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
