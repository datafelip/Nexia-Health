package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/felipelima/nexia-health/api/internal/services"
)

type AIHandler struct {
	service *services.AIService
}

func NewAIHandler(service *services.AIService) *AIHandler {
	return &AIHandler{
		service: service,
	}
}

// PreverDemanda recebe uma data e consulta o modelo
// de Machine Learning através do serviço Python.
func (h *AIHandler) PreverDemanda(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Body == nil {
		ErrorResponseJSON(
			w,
			http.StatusBadRequest,
			"Corpo da requisição obrigatório",
		)
		return
	}

	// Limita o tamanho do JSON recebido.
	r.Body = http.MaxBytesReader(
		w,
		r.Body,
		1024*1024,
	)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var request services.AIPrevisaoRequest

	if err := decoder.Decode(&request); err != nil {
		ErrorResponseJSON(
			w,
			http.StatusBadRequest,
			"JSON inválido",
		)
		return
	}

	// Impede múltiplos objetos JSON na mesma requisição.
	var extra interface{}

	if err := decoder.Decode(&extra); err != io.EOF {
		ErrorResponseJSON(
			w,
			http.StatusBadRequest,
			"O corpo deve conter apenas um objeto JSON",
		)
		return
	}

	// Encaminha os dados ao serviço de IA.
	previsao, err := h.service.PreverDemanda(
		r.Context(),
		request,
	)

	if err != nil {
		switch {
		case errors.Is(err, services.ErrAIIndisponivel):
			ErrorResponseJSON(
				w,
				http.StatusServiceUnavailable,
				"Serviço de Inteligência Artificial indisponível",
			)

		case errors.Is(err, services.ErrAIResposta):
			ErrorResponseJSON(
				w,
				http.StatusBadGateway,
				"Resposta inválida do serviço de Inteligência Artificial",
			)

		default:
			ErrorResponseJSON(
				w,
				http.StatusBadRequest,
				err.Error(),
			)
		}
		return
	}

	JSONResponse(
		w,
		http.StatusOK,
		previsao,
	)
}
