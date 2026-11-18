package models

import (
	"github.com/google/uuid"
	"time"
)

type Especialidade struct {
	ID        uuid.UUID `json:"id_especialidade"`
	Nome      string    `json:"nome"`
	Descricao *string   `json:"descricao,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
