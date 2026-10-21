package services

import (
	"context"
	"errors"
	"strings"

	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/felipelima/nexia-health/api/internal/repositories"
	"github.com/google/uuid"
)

var (
	ErrSolicitacaoDadosInvalidos = errors.New("dados da solicitação de exame inválidos")
	ErrSolicitacaoIDInvalido     = errors.New("ID da solicitação inválido")
)

type SolicitacaoExameService struct {
	repository *repositories.SolicitacaoExameRepository
}

func NewSolicitacaoExameService(
	repository *repositories.SolicitacaoExameRepository,
) *SolicitacaoExameService {
	return &SolicitacaoExameService{repository: repository}
}

func (s *SolicitacaoExameService) Create(
	ctx context.Context,
	solicitacao *models.SolicitacaoExame,
) error {
	if solicitacao == nil {
		return ErrSolicitacaoDadosInvalidos
	}

	if solicitacao.Status == "" {
		solicitacao.Status = models.ExameSolicitado
	}

	if err := validarSolicitacao(solicitacao); err != nil {
		return err
	}

	return s.repository.Create(ctx, solicitacao)
}

func (s *SolicitacaoExameService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.SolicitacaoExame, error) {
	if id == uuid.Nil {
		return nil, ErrSolicitacaoIDInvalido
	}

	return s.repository.GetByID(ctx, id)
}

func (s *SolicitacaoExameService) List(
	ctx context.Context,
	atendimentoID *uuid.UUID,
) ([]models.SolicitacaoExame, error) {
	if atendimentoID != nil && *atendimentoID == uuid.Nil {
		return nil, ErrSolicitacaoIDInvalido
	}

	return s.repository.List(ctx, atendimentoID)
}

func (s *SolicitacaoExameService) Update(
	ctx context.Context,
	solicitacao *models.SolicitacaoExame,
) error {
	if solicitacao == nil {
		return ErrSolicitacaoDadosInvalidos
	}

	if solicitacao.ID == uuid.Nil {
		return ErrSolicitacaoIDInvalido
	}

	if err := validarSolicitacao(solicitacao); err != nil {
		return err
	}

	return s.repository.Update(ctx, solicitacao)
}

func (s *SolicitacaoExameService) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	if id == uuid.Nil {
		return ErrSolicitacaoIDInvalido
	}

	return s.repository.Delete(ctx, id)
}

func validarSolicitacao(s *models.SolicitacaoExame) error {
	if s.IDAtendimento == uuid.Nil || s.IDExame == uuid.Nil {
		return ErrSolicitacaoDadosInvalidos
	}

	switch s.Status {
	case models.ExameSolicitado,
		models.ExameAgendado,
		models.ExameRealizado,
		models.ExameCancelado:
	default:
		return ErrSolicitacaoDadosInvalidos
	}

	if s.Resultado != nil {
		resultado := strings.TrimSpace(*s.Resultado)
		if resultado == "" {
			s.Resultado = nil
		} else {
			s.Resultado = &resultado
		}
	}

	return nil
}
