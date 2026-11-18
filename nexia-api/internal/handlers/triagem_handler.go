package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/felipelima/nexia-health/api/internal/repositories"
	"github.com/felipelima/nexia-health/api/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

func (h *TriagemHandler) Reclassificar(
	w http.ResponseWriter,
	r *http.Request,
) {
	// Obtém o ID da triagem pela URL.
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponseJSON(
			w,
			http.StatusBadRequest,
			"ID da triagem inválido",
		)
		return
	}

	// Recebe somente a nova classificação.
	var request struct {
		ClassificacaoRisco models.ClassificacaoRisco `json:"classificacao_risco"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponseJSON(
			w,
			http.StatusBadRequest,
			"JSON inválido",
		)
		return
	}

	// Solicita a reclassificação ao Service.
	triagem, err := h.service.Reclassificar(
		r.Context(),
		id,
		request.ClassificacaoRisco,
	)

	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrTriagemNotFound):
			ErrorResponseJSON(
				w,
				http.StatusNotFound,
				"triagem não encontrada",
			)

		case errors.Is(err, repositories.ErrTriagemNaoAguardando):
			ErrorResponseJSON(
				w,
				http.StatusConflict,
				"somente pacientes aguardando podem ser reclassificados",
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

	JSONResponse(w, http.StatusOK, triagem)
}

// Cancelar recebe uma justificativa e solicita o cancelamento
// da triagem, preservando os registros clínicos e o histórico.
func (h *TriagemHandler) Cancelar(
	w http.ResponseWriter,
	r *http.Request,
) {
	// 1. Valida o identificador recebido na URL.
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil || id == uuid.Nil {
		ErrorResponseJSON(
			w,
			http.StatusBadRequest,
			"ID da triagem inválido",
		)
		return
	}

	// 2. Limita o tamanho do corpo da requisição.
	// O motivo aceita até 1000 caracteres; o limite do
	// JSON protege contra requisições excessivamente grandes.
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	// 3. Recebe apenas o motivo do cancelamento.
	var request struct {
		Motivo string `json:"motivo"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		var maxBytesErr *http.MaxBytesError

		if errors.As(err, &maxBytesErr) {
			ErrorResponseJSON(
				w,
				http.StatusRequestEntityTooLarge,
				"corpo da requisição muito grande",
			)
			return
		}

		ErrorResponseJSON(
			w,
			http.StatusBadRequest,
			"JSON inválido: informe o campo motivo",
		)
		return
	}

	// Impede conteúdo JSON adicional após o objeto.
	var extra interface{}

	if err := decoder.Decode(&extra); err != io.EOF {
		ErrorResponseJSON(
			w,
			http.StatusBadRequest,
			"o corpo deve conter apenas um objeto JSON",
		)
		return
	}

	// 4. Valida a justificativa.
	request.Motivo = strings.TrimSpace(request.Motivo)

	if request.Motivo == "" {
		ErrorResponseJSON(
			w,
			http.StatusBadRequest,
			"o motivo do cancelamento é obrigatório",
		)
		return
	}

	if len([]rune(request.Motivo)) > 1000 {
		ErrorResponseJSON(
			w,
			http.StatusBadRequest,
			"o motivo do cancelamento deve ter no máximo 1000 caracteres",
		)
		return
	}

	// 5. Solicita ao serviço o cancelamento transacional.
	triagem, err := h.service.Cancelar(
		r.Context(),
		id,
		request.Motivo,
	)

	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrTriagemNotFound):
			ErrorResponseJSON(
				w,
				http.StatusNotFound,
				"triagem não encontrada",
			)

		case errors.Is(err, repositories.ErrTriagemJaCancelada):
			ErrorResponseJSON(
				w,
				http.StatusConflict,
				"esta triagem já foi cancelada",
			)

		case errors.Is(err, repositories.ErrTriagemCancelamentoNaoPermitido):
			ErrorResponseJSON(
				w,
				http.StatusConflict,
				"somente triagens com paciente aguardando podem ser canceladas",
			)

		case errors.Is(err, repositories.ErrTriagemComAtendimento):
			ErrorResponseJSON(
				w,
				http.StatusConflict,
				"não é permitido cancelar uma triagem vinculada a um atendimento",
			)

		case errors.Is(err, repositories.ErrMotivoCancelamentoObrigatorio):
			ErrorResponseJSON(
				w,
				http.StatusBadRequest,
				"o motivo do cancelamento é obrigatório",
			)

		default:
			// Não expõe detalhes internos do PostgreSQL ao cliente.
			ErrorResponseJSON(
				w,
				http.StatusInternalServerError,
				"erro interno ao cancelar triagem",
			)
		}
		return
	}

	// 6. Retorna a triagem atualizada com os dados
	// de cancelamento e seu histórico preservado.
	JSONResponse(w, http.StatusOK, triagem)
}
