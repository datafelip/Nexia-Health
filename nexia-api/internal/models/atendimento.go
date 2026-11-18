package models

import (
	"github.com/google/uuid"
	"time"
)

type Atendimento struct {
	ID              uuid.UUID  `json:"id_atendimento"`
	PacienteID      uuid.UUID  `json:"id_paciente"`
	ProfissionalID  uuid.UUID  `json:"id_profissional"`
	TriagemID       uuid.UUID  `json:"id_triagem"`
	AgendamentoID   *uuid.UUID `json:"id_agendamento,omitempty"`
	EspecialidadeID uuid.UUID  `json:"id_especialidade"`
	SetorID         uuid.UUID  `json:"id_setor"`
	DataInicio      time.Time  `json:"data_inicio"`
	DataFim         *time.Time `json:"data_fim,omitempty"`
	Diagnostico     *string    `json:"diagnostico,omitempty"`
	Observacoes     *string    `json:"observacoes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
