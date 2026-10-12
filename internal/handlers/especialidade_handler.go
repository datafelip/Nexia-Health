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

type EspecialidadeHandler struct {
	service *services.EspecialidadeService
}

func NewEspecialidadeHandler(service *services.EspecialidadeService) *EspecialidadeHandler {
	return &EspecialidadeHandler{
		service: service,
	}
}

func (h *EspecialidadeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var especialidade models.Especialidade

	if err := json.NewDecoder(r.Body).Decode(&especialidade); err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.service.Create(r.Context(), &especialidade); err != nil {
		if errors.Is(err, repositories.ErrEspecialidadeDuplicada) {
			ErrorResponseJSON(w, http.StatusConflict, err.Error())
			return
		}
		ErrorResponseJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	JSONResponse(w, http.StatusCreated, especialidade)
}

// List aceita o filtro opcional ?nome=<texto> (busca parcial, sem diferenciar maiúsculas/minúsculas).
func (h *EspecialidadeHandler) List(w http.ResponseWriter, r *http.Request) {
	nome := r.URL.Query().Get("nome")

	especialidades, err := h.service.List(r.Context(), nome)
	if err != nil {
		ErrorResponseJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, especialidades)
}

func (h *EspecialidadeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "ID da especialidade inválido")
		return
	}

	especialidade, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrEspecialidadeNaoEncontrada) {
			ErrorResponseJSON(w, http.StatusNotFound, err.Error())
			return
		}
		ErrorResponseJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, especialidade)
}

func (h *EspecialidadeHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "ID da especialidade inválido")
		return
	}

	var especialidade models.Especialidade

	if err := json.NewDecoder(r.Body).Decode(&especialidade); err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	especialidade.ID = id

	atualizada, err := h.service.Update(r.Context(), &especialidade)
	if err != nil {
		if errors.Is(err, repositories.ErrEspecialidadeNaoEncontrada) {
			ErrorResponseJSON(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, repositories.ErrEspecialidadeDuplicada) {
			ErrorResponseJSON(w, http.StatusConflict, err.Error())
			return
		}
		ErrorResponseJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, atualizada)
}

func (h *EspecialidadeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "ID da especialidade inválido")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrEspecialidadeNaoEncontrada) {
			ErrorResponseJSON(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, repositories.ErrEspecialidadeEmUso) {
			ErrorResponseJSON(w, http.StatusConflict, err.Error())
			return
		}
		ErrorResponseJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, map[string]string{
		"message": "especialidade excluída com sucesso",
	})
}
