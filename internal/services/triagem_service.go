package services

import (
	"context"
	"errors"
	"fmt"
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

func (s *TriagemService) Create(ctx context.Context, triagem *models.Triagem) error {
	if err := validarTriagem(triagem); err != nil {
		return err
	}

	if err := s.repository.Create(ctx, triagem); err != nil {
		return fmt.Errorf("erro ao criar triagem: %w", err)
	}

	return nil
}

func (s *TriagemService) List(ctx context.Context) ([]models.Triagem, error) {
	return s.repository.List(ctx)
}

func (s *TriagemService) GetByID(ctx context.Context, id uuid.UUID) (*models.Triagem, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *TriagemService) Update(ctx context.Context, triagem *models.Triagem) (*models.Triagem, error) {
	if err := validarTriagem(triagem); err != nil {
		return nil, err
	}

	return s.repository.Update(ctx, triagem)
}

func (s *TriagemService) Delete(ctx context.Context, id uuid.UUID) error {
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
