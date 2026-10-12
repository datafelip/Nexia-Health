package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/felipelima/nexia-health/api/internal/repositories"
	"github.com/google/uuid"
)

type SetorService struct {
	repository *repositories.SetorRepository
}

func NewSetorService(repository *repositories.SetorRepository) *SetorService {
	return &SetorService{repository: repository}
}

func (s *SetorService) Create(ctx context.Context, setor *models.Setor) error {
	setor.Nome = strings.TrimSpace(setor.Nome)

	if setor.Nome == "" {
		return fmt.Errorf("nome do setor é obrigatório")
	}

	return s.repository.Create(ctx, setor)
}

func (s *SetorService) GetByID(ctx context.Context, id uuid.UUID) (*models.Setor, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *SetorService) List(ctx context.Context, nome string) ([]models.Setor, error) {
	return s.repository.List(ctx, strings.TrimSpace(nome))
}

// Update substitui nome e descricao e retorna o registro completo após a atualização.
func (s *SetorService) Update(ctx context.Context, setor *models.Setor) (*models.Setor, error) {
	setor.Nome = strings.TrimSpace(setor.Nome)

	if setor.Nome == "" {
		return nil, fmt.Errorf("nome do setor é obrigatório")
	}

	if err := s.repository.Update(ctx, setor); err != nil {
		return nil, err
	}

	return s.repository.GetByID(ctx, setor.ID)
}

func (s *SetorService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repository.Delete(ctx, id)
}
