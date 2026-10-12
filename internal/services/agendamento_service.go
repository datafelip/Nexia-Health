package services

import (
	"context"
	"fmt"

	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/felipelima/nexia-health/api/internal/repositories"
	"github.com/google/uuid"
)

type AgendamentoService struct {
	repository *repositories.AgendamentoRepository
}

func NewAgendamentoService(repository *repositories.AgendamentoRepository) *AgendamentoService {
	return &AgendamentoService{repository: repository}
}

func tipoAgendamentoValido(tipo models.TipoAgendamento) bool {
	switch tipo {
	case models.TipoConsulta, models.TipoExame:
		return true
	default:
		return false
	}
}

func statusAgendamentoValido(status models.StatusAgendamento) bool {
	switch status {
	case models.Agendado, models.Confirmado, models.Cancelado, models.Realizado:
		return true
	default:
		return false
	}
}

func (s *AgendamentoService) Create(ctx context.Context, agendamento *models.Agendamento) error {
	if agendamento.PacienteID == uuid.Nil {
		return fmt.Errorf("id_paciente é obrigatório")
	}

	if !tipoAgendamentoValido(agendamento.Tipo) {
		return fmt.Errorf("tipo inválido, use Consulta ou Exame")
	}

	if agendamento.DataHora.IsZero() {
		return fmt.Errorf("data_hora é obrigatória")
	}

	
	if agendamento.Status == "" {
		agendamento.Status = models.Agendado
	} else if !statusAgendamentoValido(agendamento.Status) {
		return fmt.Errorf("status inválido")
	}

	return s.repository.Create(ctx, agendamento)
}

func (s *AgendamentoService) GetByID(ctx context.Context, id uuid.UUID) (*models.Agendamento, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *AgendamentoService) List(ctx context.Context, filter repositories.AgendamentoFilter) ([]models.Agendamento, error) {
	if filter.DataInicio != nil && filter.DataFim != nil && filter.DataFim.Before(*filter.DataInicio) {
		return nil, fmt.Errorf("data_fim do filtro não pode ser anterior a data_inicio")
	}

	if filter.Status != nil && !statusAgendamentoValido(*filter.Status) {
		return nil, fmt.Errorf("status de filtro inválido")
	}

	return s.repository.List(ctx, filter)
}

func (s *AgendamentoService) Update(ctx context.Context, agendamento *models.Agendamento) (*models.Agendamento, error) {
	if !tipoAgendamentoValido(agendamento.Tipo) {
		return nil, fmt.Errorf("tipo inválido, use Consulta ou Exame")
	}

	if agendamento.DataHora.IsZero() {
		return nil, fmt.Errorf("data_hora é obrigatória")
	}

	if !statusAgendamentoValido(agendamento.Status) {
		return nil, fmt.Errorf("status inválido")
	}

	if err := s.repository.Update(ctx, agendamento); err != nil {
		return nil, err
	}

	return s.repository.GetByID(ctx, agendamento.ID)
}

func (s *AgendamentoService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repository.Delete(ctx, id)
}
