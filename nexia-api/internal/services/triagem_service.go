package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/felipelima/nexia-health/api/internal/repositories"
	"github.com/google/uuid"
)

type TriagemService struct {
	repository *repositories.TriagemRepository
}

func NewTriagemService(repository *repositories.TriagemRepository) *TriagemService {
	return &TriagemService{
		repository: repository,
	}
}

func (s *TriagemService) Create(
	ctx context.Context,
	triagem *models.Triagem,
) error {
	if triagem == nil {
		return errors.New("dados da triagem são obrigatórios")
	}

	if err := validarTriagem(triagem); err != nil {
		return err
	}

	if err := s.repository.CreateComFila(ctx, triagem); err != nil {
		return fmt.Errorf(
			"erro ao criar triagem e inserir na fila: %w",
			err,
		)
	}

	return nil
}

func (s *TriagemService) List(
	ctx context.Context,
) ([]models.Triagem, error) {
	return s.repository.List(ctx)
}

func (s *TriagemService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Triagem, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *TriagemService) Update(
	ctx context.Context,
	triagem *models.Triagem,
) (*models.Triagem, error) {
	if triagem == nil {
		return nil, errors.New("dados da triagem são obrigatórios")
	}

	if err := validarTriagem(triagem); err != nil {
		return nil, err
	}

	return s.repository.Update(ctx, triagem)
}

func (s *TriagemService) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	return s.repository.Delete(ctx, id)
}

func validarTriagem(triagem *models.Triagem) error {
	if triagem.PacienteID == uuid.Nil {
		return errors.New("paciente é obrigatório")
	}

	if triagem.ProfissionalID == uuid.Nil {
		return errors.New("profissional é obrigatório")
	}

	switch triagem.ClassificacaoRisco {
	case models.Critico,
		models.Medio,
		models.Baixo:
	default:
		return errors.New("classificação de risco inválida")
	}

	return nil
}

func (s *TriagemService) Reclassificar(
	ctx context.Context,
	id uuid.UUID,
	novaClassificacao models.ClassificacaoRisco,
) (*models.Triagem, error) {
	// Verifica se o ID informado é válido.
	if id == uuid.Nil {
		return nil, errors.New("ID da triagem é obrigatório")
	}

	// Valida as classificações aceitas pelo sistema.
	switch novaClassificacao {
	case models.Critico, models.Medio, models.Baixo:
		// Classificação válida.
	default:
		return nil, errors.New(
			"classificação de risco inválida: use Vermelho, Amarelo ou Verde",
		)
	}

	// Executa a reclassificação e a reorganização da fila.
	triagem, err := s.repository.ReclassificarComFila(
		ctx,
		id,
		novaClassificacao,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"erro ao reclassificar triagem: %w",
			err,
		)
	}

	return triagem, nil
}

// Cancelar valida a solicitação de cancelamento e delega
// ao repositório a operação transacional.
//
// Regras adicionais, como impedir o cancelamento de
// triagens já atendidas ou fora da fila de espera,
// são verificadas pelo repositório.
func (s *TriagemService) Cancelar(
	ctx context.Context,
	id uuid.UUID,
	motivo string,
) (*models.Triagem, error) {
	// Impede solicitações sem um identificador válido.
	if id == uuid.Nil {
		return nil, errors.New("ID da triagem é obrigatório")
	}

	// Remove espaços desnecessários da justificativa.
	motivo = strings.TrimSpace(motivo)

	// Toda operação de cancelamento exige justificativa.
	if motivo == "" {
		return nil, repositories.ErrMotivoCancelamentoObrigatorio
	}

	// Evita justificativas excessivamente longas.
	// O limite é medido em caracteres Unicode.
	if len([]rune(motivo)) > 1000 {
		return nil, errors.New(
			"o motivo do cancelamento deve ter no máximo 1000 caracteres",
		)
	}

	// Cancela a triagem e atualiza a fila na mesma transação.
	triagem, err := s.repository.CancelarComFila(
		ctx,
		id,
		motivo,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"erro ao cancelar triagem: %w",
			err,
		)
	}

	return triagem, nil
}
