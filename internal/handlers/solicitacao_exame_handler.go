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

type SolicitacaoExameHandler struct {
	service *services.SolicitacaoExameService
}

func NewSolicitacaoExameHandler(
	service *services.SolicitacaoExameService,
) *SolicitacaoExameHandler {
	return &SolicitacaoExameHandler{service: service}
}

func (h *SolicitacaoExameHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var solicitacao models.SolicitacaoExame

	if err := json.NewDecoder(r.Body).Decode(&solicitacao); err != nil {
		writeSolicitacaoJSON(w, http.StatusBadRequest, map[string]string{
			"erro": "JSON inválido",
		})
		return
	}

	if err := h.service.Create(r.Context(), &solicitacao); err != nil {
		writeSolicitacaoError(w, err)
		return
	}

	writeSolicitacaoJSON(w, http.StatusCreated, solicitacao)
}

func (h *SolicitacaoExameHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {
	var atendimentoID *uuid.UUID

	if valor := r.URL.Query().Get("id_atendimento"); valor != "" {
		id, err := uuid.Parse(valor)
		if err != nil {
			writeSolicitacaoJSON(w, http.StatusBadRequest, map[string]string{
				"erro": "ID do atendimento inválido",
			})
			return
		}
		atendimentoID = &id
	}

	solicitacoes, err := h.service.List(r.Context(), atendimentoID)
	if err != nil {
		writeSolicitacaoError(w, err)
		return
	}

	writeSolicitacaoJSON(w, http.StatusOK, solicitacoes)
}

func (h *SolicitacaoExameHandler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeSolicitacaoJSON(w, http.StatusBadRequest, map[string]string{
			"erro": "ID da solicitação inválido",
		})
		return
	}

	solicitacao, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeSolicitacaoError(w, err)
		return
	}

	writeSolicitacaoJSON(w, http.StatusOK, solicitacao)
}

func (h *SolicitacaoExameHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeSolicitacaoJSON(w, http.StatusBadRequest, map[string]string{
			"erro": "ID da solicitação inválido",
		})
		return
	}

	var solicitacao models.SolicitacaoExame

	if err := json.NewDecoder(r.Body).Decode(&solicitacao); err != nil {
		writeSolicitacaoJSON(w, http.StatusBadRequest, map[string]string{
			"erro": "JSON inválido",
		})
		return
	}

	solicitacao.ID = id

	if err := h.service.Update(r.Context(), &solicitacao); err != nil {
		writeSolicitacaoError(w, err)
		return
	}

	writeSolicitacaoJSON(w, http.StatusOK, solicitacao)
}

func (h *SolicitacaoExameHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeSolicitacaoJSON(w, http.StatusBadRequest, map[string]string{
			"erro": "ID da solicitação inválido",
		})
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeSolicitacaoError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeSolicitacaoError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	mensagem := "Erro interno ao processar a solicitação"

	switch {
	case errors.Is(err, services.ErrSolicitacaoDadosInvalidos),
		errors.Is(err, services.ErrSolicitacaoIDInvalido),
		errors.Is(err, repositories.ErrSolicitacaoReferencia):
		status = http.StatusBadRequest
		mensagem = err.Error()

	case errors.Is(err, repositories.ErrSolicitacaoNaoEncontrada):
		status = http.StatusNotFound
		mensagem = err.Error()

	case errors.Is(err, repositories.ErrSolicitacaoDuplicada),
		errors.Is(err, repositories.ErrSolicitacaoEmUso):
		status = http.StatusConflict
		mensagem = err.Error()
	}

	writeSolicitacaoJSON(w, status, map[string]string{
		"erro": mensagem,
	})
}

func writeSolicitacaoJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
