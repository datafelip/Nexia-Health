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

type UsuarioHandler struct {
	service *services.UsuarioService
}

func NewUsuarioHandler(service *services.UsuarioService) *UsuarioHandler {
	return &UsuarioHandler{
		service: service,
	}
}

func (h *UsuarioHandler) Create(w http.ResponseWriter, r *http.Request) {
	var usuario models.Usuario

	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.service.Create(r.Context(), &usuario); err != nil {
		if errors.Is(err, repositories.ErrUsuarioDuplicado) {
			ErrorResponseJSON(w, http.StatusConflict, err.Error())
			return
		}

		ErrorResponseJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	JSONResponse(w, http.StatusCreated, usuario)
}

func (h *UsuarioHandler) List(w http.ResponseWriter, r *http.Request) {
	nome := r.URL.Query().Get("nome")

	usuarios, err := h.service.List(r.Context(), nome)
	if err != nil {
		ErrorResponseJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, usuarios)
}

func (h *UsuarioHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "ID do usuário inválido")
		return
	}

	usuario, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrUsuarioNaoEncontrado) {
			ErrorResponseJSON(w, http.StatusNotFound, err.Error())
			return
		}

		ErrorResponseJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, usuario)
}

func (h *UsuarioHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "ID do usuário inválido")
		return
	}

	var usuario models.Usuario

	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	usuario.ID = id

	atualizado, err := h.service.Update(r.Context(), &usuario)
	if err != nil {
		if errors.Is(err, repositories.ErrUsuarioNaoEncontrado) {
			ErrorResponseJSON(w, http.StatusNotFound, err.Error())
			return
		}

		if errors.Is(err, repositories.ErrUsuarioDuplicado) {
			ErrorResponseJSON(w, http.StatusConflict, err.Error())
			return
		}

		ErrorResponseJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, atualizado)
}

func (h *UsuarioHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "ID do usuário inválido")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrUsuarioNaoEncontrado) {
			ErrorResponseJSON(w, http.StatusNotFound, err.Error())
			return
		}

		ErrorResponseJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, map[string]string{
		"message": "usuário excluído com sucesso",
	})
}
