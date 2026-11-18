package models

import (
	"github.com/google/uuid"
	"time"
)

type StatusExame string

const (
	ExameSolicitado StatusExame = "Solicitado"
	ExameAgendado   StatusExame = "Agendado"
	ExameRealizado  StatusExame = "Realizado"
	ExameCancelado  StatusExame = "Cancelado"
)

type SolicitacaoExame struct {
	ID              uuid.UUID   `json:"id_solicitacao"`
	IDAtendimento   uuid.UUID   `json:"id_atendimento"`
	IDExame         uuid.UUID   `json:"id_exame"`
	DataSolicitacao time.Time   `json:"data_solicitacao"`
	Status          StatusExame `json:"status"`
	Resultado       *string     `json:"resultado,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}
