package repositories


import (
	"context"
	"testing"
	"time"
	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/google/uuid"
)

func TestTriagemRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repository := NewTriagemRepository(db)

	ctx := context.Background()

	// IDs existentes no banco para teste
	pacienteID := uuid.MustParse("3eb1ee54-4418-4dcf-8b47-0d6b2c07e5be")
	profissionalID := uuid.MustParse("4fb4f364-219c-4494-b230-088105c850ed")

	pressao := "120/80"
	temperatura := 36.5
	frequenciaCardiaca := int16(78)
	observacoes := "Teste automatizado da triagem"

	triagem := &models.Triagem{
		PacienteID:         pacienteID,
		ProfissionalID:     profissionalID,
		ClassificacaoRisco: models.Medio,
		PressaoArterial:    &pressao,
		Temperatura:        &temperatura,
		FrequenciaCardiaca: &frequenciaCardiaca,
		Observacoes:        &observacoes,
		DataTriagem:        time.Now(),
	}

	err := repository.Create(ctx, triagem)
	if err != nil {
		t.Fatalf("erro ao criar triagem: %v", err)
	}

	if triagem.ID == uuid.Nil {
		t.Fatal("ID da triagem não foi gerado")
	}

	t.Logf("Triagem criada com sucesso: %s", triagem.ID)
}