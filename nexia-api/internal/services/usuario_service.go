package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/felipelima/nexia-health/api/internal/repositories"
	"github.com/google/uuid"
)

type UsuarioService struct {
	repository *repositories.UsuarioRepository
}

func NewUsuarioService(repository *repositories.UsuarioRepository) *UsuarioService {
	return &UsuarioService{
		repository: repository,
	}
}

func (s *UsuarioService) Create(
	ctx context.Context,
	usuario *models.Usuario,
) error {
	usuario.Nome = strings.TrimSpace(usuario.Nome)
	usuario.Email = strings.TrimSpace(usuario.Email)

	if usuario.Nome == "" {
		return fmt.Errorf("nome do usuário é obrigatório")
	}

	if usuario.Email == "" {
		return fmt.Errorf("e-mail do usuário é obrigatório")
	}

	if usuario.SenhaHash == "" {
		return fmt.Errorf("senha do usuário é obrigatória")
	}

	if usuario.Perfil == "" {
		return fmt.Errorf("perfil do usuário é obrigatório")
	}

	if usuario.Status == "" {
		usuario.Status = models.Ativo
	}

	return s.repository.Create(ctx, usuario)
}

func (s *UsuarioService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Usuario, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *UsuarioService) List(
	ctx context.Context,
	nome string,
) ([]models.Usuario, error) {
	return s.repository.List(
		ctx,
		strings.TrimSpace(nome),
	)
}

func (s *UsuarioService) Update(
	ctx context.Context,
	usuario *models.Usuario,
) (*models.Usuario, error) {
	usuario.Nome = strings.TrimSpace(usuario.Nome)
	usuario.Email = strings.TrimSpace(usuario.Email)

	if usuario.Nome == "" {
		return nil, fmt.Errorf("nome do usuário é obrigatório")
	}

	if usuario.Email == "" {
		return nil, fmt.Errorf("e-mail do usuário é obrigatório")
	}

	if usuario.SenhaHash == "" {
		return nil, fmt.Errorf("senha do usuário é obrigatória")
	}

	if usuario.Perfil == "" {
		return nil, fmt.Errorf("perfil do usuário é obrigatório")
	}

	if usuario.Status == "" {
		usuario.Status = models.Ativo
	}

	if err := s.repository.Update(ctx, usuario); err != nil {
		return nil, err
	}

	return s.repository.GetByID(ctx, usuario.ID)
}

func (s *UsuarioService) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	return s.repository.Delete(ctx, id)
}
