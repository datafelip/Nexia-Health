package handlers

import (
"errors"
"encoding/json"
"github.com/felipelima/nexia-health/api/internal/models"
"github.com/felipelima/nexia-health/api/internal/repositories"
"github.com/felipelima/nexia-health/api/internal/services"
"github.com/google/uuid"
"net/http"
"time"
"github.com/go-chi/chi/v5"
)


type AtendimentoHandler struct {
service *services.AtendimentoService
}

func NewAtendimentoHandler(service *services.AtendimentoService) *AtendimentoHandler {
return &AtendimentoHandler{
service: service,
    }
}

func (h *AtendimentoHandler) Create(w http.ResponseWriter, r *http.Request) {
var atendimento models.Atendimento

if err := json.NewDecoder(r.Body).Decode(&atendimento); err != nil {
ErrorResponseJSON(w, http.StatusBadRequest, "JSON inválido")
return
    }

if err := h.service.Create(r.Context(), &atendimento); err != nil {
ErrorResponseJSON(w, http.StatusBadRequest, err.Error())
return
    }

JSONResponse(w, http.StatusCreated, atendimento)
}

// List aceita filtros opcionais via query string:
//   ?id_paciente=<uuid>
//   ?data_inicio=2026-09-01
//   ?data_fim=2026-09-30
// Datas são aceitas no formato YYYY-MM-DD (data_fim é tratado como
// o fim do dia, 23:59:59, para incluir atendimentos daquele dia).
func (h *AtendimentoHandler) List(w http.ResponseWriter, r *http.Request) {
var filter repositories.AtendimentoFilter

if pacienteParam := r.URL.Query().Get("id_paciente"); pacienteParam != "" {
pacienteID, err := uuid.Parse(pacienteParam)
if err != nil {
ErrorResponseJSON(w, http.StatusBadRequest, "id_paciente inválido")
return
        }
filter.PacienteID = &pacienteID
    }

if dataInicioParam := r.URL.Query().Get("data_inicio"); dataInicioParam != "" {
dataInicio, err := time.Parse("2006-01-02", dataInicioParam)
if err != nil {
ErrorResponseJSON(w, http.StatusBadRequest, "data_inicio inválida, use o formato AAAA-MM-DD")
return
        }
filter.DataInicio = &dataInicio
    }

if dataFimParam := r.URL.Query().Get("data_fim"); dataFimParam != "" {
dataFim, err := time.Parse("2006-01-02", dataFimParam)
if err != nil {
ErrorResponseJSON(w, http.StatusBadRequest, "data_fim inválida, use o formato AAAA-MM-DD")
return
        }
dataFim = dataFim.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
filter.DataFim = &dataFim
    }

atendimentos, err := h.service.List(r.Context(), filter)
if err != nil {
ErrorResponseJSON(w, http.StatusBadRequest, err.Error())
return
    }

JSONResponse(w, http.StatusOK, atendimentos)
}

func (h *AtendimentoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
idParam := chi.URLParam(r, "id")

id, err := uuid.Parse(idParam)
if err != nil {
ErrorResponseJSON(w, http.StatusBadRequest, "ID do atendimento inválido")
return
    }

atendimento, err := h.service.GetByID(r.Context(), id)
if err != nil {
if errors.Is(err, repositories.ErrAtendimentoNaoEncontrado) {
ErrorResponseJSON(w, http.StatusNotFound, err.Error())
return
        }
ErrorResponseJSON(w, http.StatusInternalServerError, err.Error())
return
    }

JSONResponse(w, http.StatusOK, atendimento)
}

func (h *AtendimentoHandler) Update(w http.ResponseWriter, r *http.Request) {
idParam := chi.URLParam(r, "id")

id, err := uuid.Parse(idParam)
if err != nil {
ErrorResponseJSON(w, http.StatusBadRequest, "ID do atendimento inválido")
return
    }

var atendimento models.Atendimento

if err := json.NewDecoder(r.Body).Decode(&atendimento); err != nil {
ErrorResponseJSON(w, http.StatusBadRequest, "JSON inválido")
return
    }

atendimento.ID = id

atualizado, err := h.service.Update(r.Context(), &atendimento)
if err != nil {
if errors.Is(err, repositories.ErrAtendimentoNaoEncontrado) {
ErrorResponseJSON(w, http.StatusNotFound, err.Error())
return
        }
ErrorResponseJSON(w, http.StatusBadRequest, err.Error())
return
    }

JSONResponse(w, http.StatusOK, atualizado)
}

func (h *AtendimentoHandler) Delete(w http.ResponseWriter, r *http.Request) {
idParam := chi.URLParam(r, "id")

id, err := uuid.Parse(idParam)
if err != nil {
ErrorResponseJSON(w, http.StatusBadRequest, "ID do atendimento inválido")
return
    }

if err := h.service.Delete(r.Context(), id); err != nil {
if errors.Is(err, repositories.ErrAtendimentoNaoEncontrado) {
ErrorResponseJSON(w, http.StatusNotFound, err.Error())
return
        }
ErrorResponseJSON(w, http.StatusInternalServerError, err.Error())
return
    }

JSONResponse(w, http.StatusOK, map[string]string{
"message": "atendimento excluído com sucesso",
    })
}