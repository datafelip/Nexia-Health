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


type AgendamentoHandler struct {
service *services.AgendamentoService
}

func NewAgendamentoHandler(service *services.AgendamentoService) *AgendamentoHandler {
return &AgendamentoHandler{
service: service,
    }
}

func (h *AgendamentoHandler) Create(w http.ResponseWriter, r *http.Request) {
var agendamento models.Agendamento

if err := json.NewDecoder(r.Body).Decode(&agendamento); err != nil {
ErrorResponseJSON(w, http.StatusBadRequest, "JSON inválido")
return
    }

if err := h.service.Create(r.Context(), &agendamento); err != nil {
ErrorResponseJSON(w, http.StatusBadRequest, err.Error())
return
    }

JSONResponse(w, http.StatusCreated, agendamento)
}

// List aceita filtros opcionais via query string:
//   ?id_paciente=<uuid>
//   ?data_inicio=2026-09-01
//   ?data_fim=2026-09-30
//   ?status=Agendado
func (h *AgendamentoHandler) List(w http.ResponseWriter, r *http.Request) {
var filter repositories.AgendamentoFilter

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
filter.DataFim = &dataFim
    }

if statusParam := r.URL.Query().Get("status"); statusParam != "" {
status := models.StatusAgendamento(statusParam)
filter.Status = &status
    }

agendamentos, err := h.service.List(r.Context(), filter)
if err != nil {
ErrorResponseJSON(w, http.StatusBadRequest, err.Error())
return
    }

JSONResponse(w, http.StatusOK, agendamentos)
}

func (h *AgendamentoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
idParam := chi.URLParam(r, "id")

id, err := uuid.Parse(idParam)
if err != nil {
ErrorResponseJSON(w, http.StatusBadRequest, "ID do agendamento inválido")
return
    }

agendamento, err := h.service.GetByID(r.Context(), id)
if err != nil {
if errors.Is(err, repositories.ErrAgendamentoNaoEncontrado) {
ErrorResponseJSON(w, http.StatusNotFound, err.Error())
return
        }
ErrorResponseJSON(w, http.StatusInternalServerError, err.Error())
return
    }

JSONResponse(w, http.StatusOK, agendamento)
}

func (h *AgendamentoHandler) Update(w http.ResponseWriter, r *http.Request) {
idParam := chi.URLParam(r, "id")

id, err := uuid.Parse(idParam)
if err != nil {
ErrorResponseJSON(w, http.StatusBadRequest, "ID do agendamento inválido")
return
    }

var agendamento models.Agendamento

if err := json.NewDecoder(r.Body).Decode(&agendamento); err != nil {
ErrorResponseJSON(w, http.StatusBadRequest, "JSON inválido")
return
    }

agendamento.ID = id

atualizado, err := h.service.Update(r.Context(), &agendamento)
if err != nil {
if errors.Is(err, repositories.ErrAgendamentoNaoEncontrado) {
ErrorResponseJSON(w, http.StatusNotFound, err.Error())
return
        }
ErrorResponseJSON(w, http.StatusBadRequest, err.Error())
return
    }

JSONResponse(w, http.StatusOK, atualizado)
}

func (h *AgendamentoHandler) Delete(w http.ResponseWriter, r *http.Request) {
idParam := chi.URLParam(r, "id")

id, err := uuid.Parse(idParam)
if err != nil {
ErrorResponseJSON(w, http.StatusBadRequest, "ID do agendamento inválido")
return
    }

if err := h.service.Delete(r.Context(), id); err != nil {
if errors.Is(err, repositories.ErrAgendamentoNaoEncontrado) {
ErrorResponseJSON(w, http.StatusNotFound, err.Error())
return
        }
ErrorResponseJSON(w, http.StatusInternalServerError, err.Error())
return
    }

JSONResponse(w, http.StatusOK, map[string]string{
"message": "agendamento excluído com sucesso",
    })
}