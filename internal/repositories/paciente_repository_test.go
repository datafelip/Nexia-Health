package repositories

import (
	"context"
	"testing"
	"time"
	"fmt"
	"github.com/google/uuid"
	"github.com/felipelima/nexia-health/api/configs"
	"github.com/felipelima/nexia-health/api/internal/database"
	"github.com/felipelima/nexia-health/api/internal/models"
)

func TestPacienteRepository_Create(t *testing.T) {
	config, err := configs.Load()
	if err != nil {
		t.Fatalf("erro ao carregar configuração: %v", err)
	}

	db, err := database.Connect(config.DatabaseURL)
	if err != nil {
		t.Fatalf("erro ao conectar ao banco: %v", err)
	}
	defer db.Close()

	repository := NewPacienteRepository(db)

	telefone := fmt.Sprintf("1699%07d", time.Now().UnixNano()%10000000)
	endereco := "Rua Teste, 123"
	cartaoSUS := fmt.Sprintf("%015d", time.Now().UnixNano()%1000000000000000)

	paciente := &models.Paciente{
		ID:             uuid.Nil,
		Nome:           "Paciente Teste",
		CPF: fmt.Sprintf("%011d", time.Now().UnixNano()%100000000000),
		DataNascimento: time.Date(2000, 1, 15, 0, 0, 0, 0, time.UTC),
		Sexo:           models.SexoMasculino,
		Telefone:       &telefone,
		Endereco:       &endereco,
		CartaoSUS:      &cartaoSUS,
	}

	err = repository.Create(context.Background(), paciente)
	if err != nil {
		t.Fatalf("erro ao criar paciente: %v", err)
	}

	if paciente.ID == uuid.Nil {
		t.Fatal("UUID do paciente não foi gerado")
	}

	t.Logf("Paciente criado com sucesso: %s", paciente.ID)
}