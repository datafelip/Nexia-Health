package models

import (
	"time"
	"github.com/google/uuid"
)

type Sexo string

const (
	SexoMasculino Sexo = "Masculino"
	SexoFeminino Sexo = "Feminino"
	SexoOutro     Sexo = "Outro"
)

type Paciente struct {
	ID             uuid.UUID `json:"id" db:"id_paciente"`
	Nome           string    `json:"nome" db:"nome"`
	CPF            string    `json:"cpf" db:"cpf"`
	DataNascimento time.Time `json:"data_nascimento" db:"data_nascimento"`
	Sexo           Sexo      `json:"sexo" db:"sexo"`
	Telefone       *string   `json:"telefone,omitempty" db:"telefone"`
	Endereco       *string   `json:"endereco,omitempty" db:"endereco"`
	CartaoSUS      *string   `json:"cartao_sus,omitempty" db:"cartao_sus"`
	DataCadastro   time.Time `json:"data_cadastro" db:"data_cadastro"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}