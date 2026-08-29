package models

import (
	"github.com/google/uuid"
	"time"
)

type ClassificacaoRisco string

const (
	Critico ClassificacaoRisco = "Vermelho"
	Medio   ClassificacaoRisco = "Amarelo"
	Baixo   ClassificacaoRisco = "Verde"
)

type Triagem struct {
	ID                 uuid.UUID          `json:"id_triagem"`
	PacienteID         uuid.UUID          `json:"id_paciente"`
	ProfissionalID     uuid.UUID          `json:"id_profissional"`
	ClassificacaoRisco ClassificacaoRisco `json:"classificacao_risco"`
	PressaoArterial    *string            `json:"pressao_arterial,omitempty"`
	Temperatura        *float64           `json:"temperatura,omitempty"`
	FrequenciaCardiaca *int16             `json:"frequencia_cardiaca,omitempty"`
	Observacoes        *string            `json:"observacoes,omitempty"`
	DataTriagem        time.Time          `json:"data_triagem"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
}
