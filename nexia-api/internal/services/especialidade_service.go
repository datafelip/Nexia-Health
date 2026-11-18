package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/felipelima/nexia-health/api/internal/repositories"
	"github.com/google/uuid"
)

type EspecialidadeService struct {
	repository *repositories.EspecialidadeRepository
}

func NewEspecialidadeService(repository *repositories.EspecialidadeRepository) *EspecialidadeService {
	return &EspecialidadeService{repository: repository}
}

func (s *EspecialidadeService) Create(ctx context.Context, especialidade *models.Especialidade) error {
	especialidade.Nome = strings.TrimSpace(especialidade.Nome)

	if especialidade.Nome == "" {
		return fmt.Errorf("nome da especialidade é obrigatório")
	}

	return s.repository.Create(ctx, especialidade)
}

func (s *EspecialidadeService) GetByID(ctx context.Context, id uuid.UUID) (*models.Especialidade, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *EspecialidadeService) List(ctx context.Context, nome string) ([]models.Especialidade, error) {
	return s.repository.List(ctx, strings.TrimSpace(nome))
}

// Update substitui nome e descricao e retorna o registro completo após a atualização.
func (s *EspecialidadeService) Update(ctx context.Context, especialidade *models.Especialidade) (*models.Especialidade, error) {
	especialidade.Nome = strings.TrimSpace(especialidade.Nome)

	if especialidade.Nome == "" {
		return nil, fmt.Errorf("nome da especialidade é obrigatório")
	}

	if err := s.repository.Update(ctx, especialidade); err != nil {
		return nil, err
	}

	return s.repository.GetByID(ctx, especialidade.ID)
}

func (s *EspecialidadeService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repository.Delete(ctx, id)
}
