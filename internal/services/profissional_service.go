package services

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/felipelima/nexia-health/api/internal/repositories"
	"github.com/google/uuid"
)

var (
	ErrProfissionalDadosInvalidos = errors.New("dados do profissional inválidos")
	ErrProfissionalIDInvalido     = errors.New("ID do profissional inválido")
)

var cpfRegex = regexp.MustCompile(`^[0-9]{11}$`)

type ProfissionalService struct {
	repository *repositories.ProfissionalRepository
}

func NewProfissionalService(
	repository *repositories.ProfissionalRepository,
) *ProfissionalService {
	return &ProfissionalService{repository: repository}
}

func (s *ProfissionalService) Create(
	ctx context.Context,
	p *models.Profissional,
) error {
	if err := validarProfissional(p, false); err != nil {
		return err
	}

	if p.Status == "" {
		p.Status = models.ProfissionalAtivo
	}

	return s.repository.Create(ctx, p)
}

func (s *ProfissionalService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Profissional, error) {
	if id == uuid.Nil {
		return nil, ErrProfissionalIDInvalido
	}

	return s.repository.GetByID(ctx, id)
}

func (s *ProfissionalService) List(
	ctx context.Context,
	cpf string,
) ([]models.Profissional, error) {
	cpf = strings.TrimSpace(cpf)

	if cpf != "" && !regexp.MustCompile(`^[0-9]{1,11}$`).MatchString(cpf) {
		return nil, ErrProfissionalDadosInvalidos
	}

	return s.repository.List(ctx, cpf)
}

func (s *ProfissionalService) Update(
	ctx context.Context,
	p *models.Profissional,
) error {
	if p.ID == uuid.Nil {
		return ErrProfissionalIDInvalido
	}

	if err := validarProfissional(p, true); err != nil {
		return err
	}

	return s.repository.Update(ctx, p)
}

func (s *ProfissionalService) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	if id == uuid.Nil {
		return ErrProfissionalIDInvalido
	}

	return s.repository.Delete(ctx, id)
}

func validarProfissional(p *models.Profissional, exigirStatus bool) error {
	if p == nil ||
		p.IDUsuario == uuid.Nil ||
		p.IDEspecialidade == uuid.Nil ||
		p.IDSetor == uuid.Nil ||
		!cpfRegex.MatchString(strings.TrimSpace(p.CPF)) ||
		strings.TrimSpace(p.ConselhoProfissional) == "" ||
		strings.TrimSpace(p.NumeroRegistro) == "" {
		return ErrProfissionalDadosInvalidos
	}

	p.CPF = strings.TrimSpace(p.CPF)
	p.ConselhoProfissional = strings.TrimSpace(p.ConselhoProfissional)
	p.NumeroRegistro = strings.TrimSpace(p.NumeroRegistro)

	if p.Telefone != nil {
		telefone := strings.TrimSpace(*p.Telefone)
		if telefone == "" {
			p.Telefone = nil
		} else {
			p.Telefone = &telefone
		}
	}

	if p.Status == "" && !exigirStatus {
		return nil
	}

	switch p.Status {
	case models.ProfissionalAtivo,
		models.ProfissionalInativo,
		models.ProfissionalBloqueado:
		return nil
	default:
		return ErrProfissionalDadosInvalidos
	}
}
