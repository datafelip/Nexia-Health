package models

import (
	"github.com/google/uuid"
	"time"
)

type Status string

const (
	Aguardando    Status = "Aguardando"
	Chamado       Status = "Chamado"
	EmAtendimento Status = "Em Atendimento"
	Finalizado    Status = "Finalizado"
	Canceled    Status = "Cancelado"
)

type Fila struct {
	ID          uuid.UUID `json:"id_fila"`
	PacienteID  uuid.UUID `json:"id_paciente"`
	TriagemID   uuid.UUID `json:"id_triagem"`
	Posicao     int       `json:"posicao"`
	Status      Status    `json:"status"`
	DataEntrada time.Time `json:"data_entrada"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
