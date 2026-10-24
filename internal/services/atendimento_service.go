package services

import (
	"context"
	"fmt"
	"time"

	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/felipelima/nexia-health/api/internal/repositories"
	"github.com/google/uuid"
)

type AtendimentoService struct {
	repository *repositories.AtendimentoRepository
}

func NewAtendimentoService(repository *repositories.AtendimentoRepository) *AtendimentoService {
	return &AtendimentoService{repository: repository}
}

func (s *AtendimentoService) Create(
	ctx context.Context,
	atendimento *models.Atendimento,
) error {
	if atendimento == nil {
		return fmt.Errorf("dados do atendimento são obrigatórios")
	}

	if atendimento.PacienteID == uuid.Nil {
		return fmt.Errorf("id_paciente é obrigatório")
	}

	if atendimento.ProfissionalID == uuid.Nil {
		return fmt.Errorf("id_profissional é obrigatório")
	}

	if atendimento.TriagemID == uuid.Nil {
		return fmt.Errorf("id_triagem é obrigatório")
	}

	if atendimento.EspecialidadeID == uuid.Nil {
		return fmt.Errorf("id_especialidade é obrigatório")
	}

	if atendimento.SetorID == uuid.Nil {
		return fmt.Errorf("id_setor é obrigatório")
	}

	if atendimento.DataInicio.IsZero() {
		atendimento.DataInicio = time.Now()
	}

	if atendimento.DataFim != nil {
		return fmt.Errorf(
			"um novo atendimento deve ser criado sem data_fim",
		)
	}

	if err := s.repository.CreateComFila(
		ctx, atendimento,
	); err != nil {
		return fmt.Errorf(
			"erro ao iniciar atendimento: %w", err,
		)
	}

	return nil
}

func (s *AtendimentoService) GetByID(ctx context.Context, id uuid.UUID) (*models.Atendimento, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *AtendimentoService) List(ctx context.Context, filter repositories.AtendimentoFilter) ([]models.Atendimento, error) {
	if filter.DataInicio != nil && filter.DataFim != nil && filter.DataFim.Before(*filter.DataInicio) {
		return nil, fmt.Errorf("data_fim do filtro não pode ser anterior a data_inicio")
	}

	return s.repository.List(ctx, filter)
}

func (s *AtendimentoService) Update(
	ctx context.Context,
	atendimento *models.Atendimento,
) (*models.Atendimento, error) {
	if atendimento == nil {
		return nil, fmt.Errorf("dados do atendimento são obrigatórios")
	}

	if atendimento.DataFim == nil {
		return nil, fmt.Errorf(
			"data_fim é obrigatória para finalizar o atendimento",
		)
	}

	if err := s.repository.FinalizarComFila(
		ctx, atendimento,
	); err != nil {
		return nil, fmt.Errorf(
			"erro ao finalizar atendimento: %w", err,
		)
	}

	return s.repository.GetByID(ctx, atendimento.ID)
}

func (s *AtendimentoService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repository.Delete(ctx, id)
}
