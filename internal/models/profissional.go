
package models

import (
	"time"

	"github.com/google/uuid"
)

type StatusProfissional string

const (
	ProfissionalAtivo     StatusProfissional = "Ativo"
	ProfissionalInativo   StatusProfissional = "Inativo"
	ProfissionalBloqueado StatusProfissional = "Bloqueado"
)

type Profissional struct {
	ID                 uuid.UUID         `json:"id_profissional"`
	IDUsuario          uuid.UUID         `json:"id_usuario"`
	IDEspecialidade    uuid.UUID         `json:"id_especialidade"`
	IDSetor             uuid.UUID         `json:"id_setor"`
	CPF                 string            `json:"cpf"`
	ConselhoProfissional string            `json:"conselho_profissional"`
	NumeroRegistro      string            `json:"numero_registro"`
	Telefone            *string           `json:"telefone,omitempty"`
	Status              StatusProfissional `json:"status"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`
}
