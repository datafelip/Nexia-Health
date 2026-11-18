package models

import (
	"time"

	"github.com/google/uuid"
)

type Exame struct {
	ID        uuid.UUID `json:"id_exame"`
	Nome      string    `json:"nome"`
	Descricao *string   `json:"descricao,omitempty"`
	Preparo   *string   `json:"preparo,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
