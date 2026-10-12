package models

import (
	"time"

	"github.com/google/uuid"
)

type StatusAgendamento string

const (
	Agendado   StatusAgendamento = "Agendado"
	Confirmado StatusAgendamento = "Confirmado"
	Cancelado  StatusAgendamento = "Cancelado"
	Realizado  StatusAgendamento = "Realizado"
)

type TipoAgendamento string

const (
	TipoConsulta TipoAgendamento = "Consulta"
	TipoExame    TipoAgendamento = "Exame"
)

type Agendamento struct {
	ID         uuid.UUID         `json:"id_agendamento"`
	PacienteID uuid.UUID         `json:"id_paciente"`
	Tipo       TipoAgendamento   `json:"tipo"`
	DataHora   time.Time         `json:"data_hora"`
	Status     StatusAgendamento `json:"status"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}
