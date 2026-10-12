package models

import (
	"time"

	"github.com/google/uuid"
)

type PerfilUsuario string

const (
	Administrador PerfilUsuario = "Administrador"
	Recepcionista PerfilUsuario = "Recepcionista"
	Enfermeiro    PerfilUsuario = "Enfermeiro"
	Medico        PerfilUsuario = "Médico"
	Gestor        PerfilUsuario = "Gestor"
)

type StatusUsuario string

const (
	Ativo     StatusUsuario = "Ativo"
	Inativo   StatusUsuario = "Inativo"
	Bloqueado StatusUsuario = "Bloqueado"
)

type Usuario struct {
	ID          uuid.UUID     `json:"id_usuario"`
	Nome        string        `json:"nome"`
	Email       string        `json:"email"`
	SenhaHash   string        `json:"senha_hash"`
	Perfil      PerfilUsuario `json:"perfil"`
	Status      StatusUsuario `json:"status"`
	UltimoLogin *time.Time    `json:"ultimo_login,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}
