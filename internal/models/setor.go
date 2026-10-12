package models

import (
	"github.com/google/uuid"
	"time"
)

type Setor struct {
	ID        uuid.UUID `json:"id_setor"`
	Nome      string    `json:"nome"`
	Descricao *string   `json:"descricao,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
