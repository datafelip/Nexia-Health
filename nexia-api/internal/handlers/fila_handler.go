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

type FilaHandler struct {
	service *services.FilaService
}

func NewFilaHandler(service *services.FilaService) *FilaHandler {
	return &FilaHandler{
		service: service,
	}
}

func (h *FilaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var fila models.Fila

	if err := json.NewDecoder(r.Body).Decode(&fila); err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.service.Create(r.Context(), &fila); err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	JSONResponse(w, http.StatusCreated, fila)
}

func (h *FilaHandler) List(w http.ResponseWriter, r *http.Request) {
	filas, err := h.service.List(r.Context())
	if err != nil {
		ErrorResponseJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, filas)
}

func (h *FilaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "ID da fila inválido")
		return
	}

	fila, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrFilaNaoEncontrada) {
			ErrorResponseJSON(w, http.StatusNotFound, err.Error())
			return
		}
		ErrorResponseJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, fila)
}

func (h *FilaHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "ID da fila inválido")
		return
	}

	var fila models.Fila

	if err := json.NewDecoder(r.Body).Decode(&fila); err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	fila.ID = id

	atualizada, err := h.service.Update(r.Context(), &fila)
	if err != nil {
		if errors.Is(err, repositories.ErrFilaNaoEncontrada) {
			ErrorResponseJSON(w, http.StatusNotFound, err.Error())
			return
		}
		ErrorResponseJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, atualizada)
}

func (h *FilaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponseJSON(w, http.StatusBadRequest, "ID da fila inválido")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrFilaNaoEncontrada) {
			ErrorResponseJSON(w, http.StatusNotFound, err.Error())
			return
		}
		ErrorResponseJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, map[string]string{
		"message": "fila excluída com sucesso",
	})
}

func (h *FilaHandler) ChamarProximo(
	w http.ResponseWriter,
	r *http.Request,
) {
	fila, err := h.service.ChamarProximo(r.Context())

	if err != nil {
		if errors.Is(err, repositories.ErrFilaVazia) {
			ErrorResponseJSON(
				w,
				http.StatusConflict,
				err.Error(),
			)
			return
		}

		ErrorResponseJSON(
			w,
			http.StatusInternalServerError,
			"erro interno ao chamar próximo paciente",
		)
		return
	}

	JSONResponse(w, http.StatusOK, fila)
}

func (h *FilaHandler) Cancelar(
	w http.ResponseWriter,
	r *http.Request,
) {
	// Obtém o ID da fila pela URL.
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil || id == uuid.Nil {
		ErrorResponseJSON(
			w,
			http.StatusBadRequest,
			"ID da fila inválido",
		)
		return
	}

	// Solicita o cancelamento ao Service.
	err = h.service.Cancelar(r.Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrFilaNaoEncontrada):
			ErrorResponseJSON(
				w,
				http.StatusNotFound,
				"fila não encontrada",
			)

		case errors.Is(err, repositories.ErrFilaNaoCancelavel):
			ErrorResponseJSON(
				w,
				http.StatusConflict,
				"somente pacientes aguardando ou chamados podem ter a fila cancelada",
			)

		default:
			ErrorResponseJSON(
				w,
				http.StatusInternalServerError,
				"erro interno ao cancelar fila",
			)
		}
		return
	}

	// Cancelamento concluído.
	w.WriteHeader(http.StatusNoContent)
}
