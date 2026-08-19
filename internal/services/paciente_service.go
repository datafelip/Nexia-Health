package services

import (
	"context"
	"fmt"
	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/felipelima/nexia-health/api/internal/repositories"
	"github.com/google/uuid"
	"strings"
)

type PacienteService struct {
	repository *repositories.PacienteRepository
}

func NewPacienteService(repository *repositories.PacienteRepository) *PacienteService {
	return &PacienteService{
		repository: repository,
	}
}

func (s *PacienteService) Create(ctx context.Context, paciente *models.Paciente) error {
	if strings.TrimSpace(paciente.Nome) == "" {
		return fmt.Errorf("nome do paciente é obrigatório")
	}

	if strings.TrimSpace(paciente.CPF) == "" {
		return fmt.Errorf("CPF do paciente é obrigatório")
	}

	if paciente.DataNascimento.IsZero() {
		return fmt.Errorf("data de nascimento é obrigatória")
	}

	if paciente.Sexo == "" {
		return fmt.Errorf("sexo do paciente é obrigatório")
	}

	if err := s.repository.Create(ctx, paciente); err != nil {
		return err
	}

	return nil
}

func (s *PacienteService) List(ctx context.Context) ([]models.Paciente, error) {
	pacientes, err := s.repository.List(ctx)
	if err != nil {
		return nil, err
	}

	return pacientes, nil
}

func (s *PacienteService) GetByID(ctx context.Context, id uuid.UUID) (*models.Paciente, error) {
	paciente, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return paciente, nil
}

func (s *PacienteService) Update(ctx context.Context, paciente *models.Paciente) (*models.Paciente, error) {
	if strings.TrimSpace(paciente.Nome) == "" {
		return nil, fmt.Errorf("nome do paciente é obrigatório")
	}

	if strings.TrimSpace(paciente.CPF) == "" {
		return nil, fmt.Errorf("CPF do paciente é obrigatório")
	}

	if paciente.DataNascimento.IsZero() {
		return nil, fmt.Errorf("data de nascimento é obrigatória")
	}

	if paciente.Sexo == "" {
		return nil, fmt.Errorf("sexo do paciente é obrigatório")
	}

	if err := s.repository.Update(ctx, paciente); err != nil {
		return nil, err
	}

	return s.repository.GetByID(ctx, paciente.ID)
}

func (s *PacienteService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("ID do paciente inválido")
	}
	
	if err := s.repository.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}