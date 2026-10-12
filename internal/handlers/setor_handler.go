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

type SetorHandler struct {
	service *services.SetorService
}

func NewSetorHandler(service *services.SetorService) *SetorHandler {
	return &SetorHandler{
		service: service,
	}
}

func (h *SetorHandler) Create(w http.ResponseWriter, r *http.Request) {
	var setor models.Setor

	if err := json.NewDecoder(r.Body).Decode(&setor); err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.service.Create(r.Context(), &setor); err != nil {
		if errors.Is(err, repositories.ErrSetorDuplicado) {
			ErrorResponseJSON(w, http.StatusConflict, err.Error())
			return
		}
		ErrorResponseJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	JSONResponse(w, http.StatusCreated, setor)
}

// List aceita o filtro opcional ?nome=<texto> (busca parcial, sem diferenciar maiúsculas/minúsculas).
func (h *SetorHandler) List(w http.ResponseWriter, r *http.Request) {
	nome := r.URL.Query().Get("nome")

	setores, err := h.service.List(r.Context(), nome)
	if err != nil {
		ErrorResponseJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, setores)
}

func (h *SetorHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "ID do setor inválido")
		return
	}

	setor, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrSetorNaoEncontrado) {
			ErrorResponseJSON(w, http.StatusNotFound, err.Error())
			return
		}
		ErrorResponseJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, setor)
}

func (h *SetorHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "ID do setor inválido")
		return
	}

	var setor models.Setor

	if err := json.NewDecoder(r.Body).Decode(&setor); err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	setor.ID = id

	atualizado, err := h.service.Update(r.Context(), &setor)
	if err != nil {
		if errors.Is(err, repositories.ErrSetorNaoEncontrado) {
			ErrorResponseJSON(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, repositories.ErrSetorDuplicado) {
			ErrorResponseJSON(w, http.StatusConflict, err.Error())
			return
		}
		ErrorResponseJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, atualizado)
}

func (h *SetorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "ID do setor inválido")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrSetorNaoEncontrado) {
			ErrorResponseJSON(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, repositories.ErrSetorEmUso) {
			ErrorResponseJSON(w, http.StatusConflict, err.Error())
			return
		}
		ErrorResponseJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, map[string]string{
		"message": "setor excluído com sucesso",
	})
}
