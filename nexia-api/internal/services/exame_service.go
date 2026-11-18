package services

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/felipelima/nexia-health/api/internal/repositories"
	"github.com/google/uuid"
)

var (
	ErrExameDadosInvalidos = errors.New("dados do exame inválidos")
	ErrExameIDInvalido     = errors.New("ID do exame inválido")
)

type ExameService struct {
	repository *repositories.ExameRepository
}

func NewExameService(
	repository *repositories.ExameRepository,
) *ExameService {
	return &ExameService{repository: repository}
}

func (s *ExameService) Create(
	ctx context.Context,
	exame *models.Exame,
) error {
	if err := validarExame(exame); err != nil {
		return err
	}

	return s.repository.Create(ctx, exame)
}

func (s *ExameService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Exame, error) {
	if id == uuid.Nil {
		return nil, ErrExameIDInvalido
	}

	return s.repository.GetByID(ctx, id)
}

func (s *ExameService) List(
	ctx context.Context,
	nome string,
) ([]models.Exame, error) {
	return s.repository.List(ctx, strings.TrimSpace(nome))
}

func (s *ExameService) Update(
	ctx context.Context,
	exame *models.Exame,
) error {
	if exame == nil {
		return ErrExameDadosInvalidos
	}

	if exame.ID == uuid.Nil {
		return ErrExameIDInvalido
	}

	if err := validarExame(exame); err != nil {
		return err
	}

	return s.repository.Update(ctx, exame)
}

func (s *ExameService) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	if id == uuid.Nil {
		return ErrExameIDInvalido
	}

	return s.repository.Delete(ctx, id)
}

func validarExame(exame *models.Exame) error {
	if exame == nil {
		return ErrExameDadosInvalidos
	}

	exame.Nome = strings.TrimSpace(exame.Nome)

	if exame.Nome == "" || utf8.RuneCountInString(exame.Nome) > 100 {
		return ErrExameDadosInvalidos
	}

	if exame.Descricao != nil {
		descricao := strings.TrimSpace(*exame.Descricao)
		exame.Descricao = &descricao
	}

	if exame.Preparo != nil {
		preparo := strings.TrimSpace(*exame.Preparo)
		exame.Preparo = &preparo
	}

	return nil
}
