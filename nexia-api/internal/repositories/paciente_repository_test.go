package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestPacienteRepository_Create(t *testing.T) {
	// setupTestDB exige TEST_DATABASE_URL apontando
	// exclusivamente para nexia_test_db.
	db := setupTestDB(t)

	repository := NewPacienteRepository(db)
	ctx := context.Background()

	// Identificador único para evitar conflitos
	// entre execuções dos testes.
	identificador := uuid.New().String()

	// CPF fictício de 11 dígitos, usado somente
	// no banco de testes.
	cpf := fmt.Sprintf(
		"%011d",
		time.Now().UnixNano()%100000000000,
	)

	telefone := "16999999999"
	endereco := "Rua Teste Automatizado, 123"

	paciente := &models.Paciente{
		Nome:           "Paciente Teste Automatizado " + identificador[:8],
		CPF:            cpf,
		DataNascimento: time.Date(2000, 1, 15, 0, 0, 0, 0, time.UTC),
		Sexo:           models.SexoMasculino,
		Telefone:       &telefone,
		Endereco:       &endereco,
	}

	// Cria o paciente no banco de testes.
	err := repository.Create(ctx, paciente)
	if err != nil {
		t.Fatalf("erro ao criar paciente: %v", err)
	}

	// Registra a limpeza imediatamente após a criação.
	// O registro será removido mesmo se uma
	// verificação posterior falhar.
	t.Cleanup(func() {
		_, err := db.Exec(
			context.Background(),
			"DELETE FROM paciente WHERE id_paciente = $1",
			paciente.ID,
		)
		if err != nil {
			t.Errorf("erro ao limpar paciente de teste: %v", err)
		}
	})

	if paciente.ID == uuid.Nil {
		t.Fatal("UUID do paciente não foi gerado")
	}

	// Confirma que o registro realmente existe.
	var nomeSalvo string

	err = db.QueryRow(
		ctx,
		"SELECT nome FROM paciente WHERE id_paciente = $1",
		paciente.ID,
	).Scan(&nomeSalvo)

	if err != nil {
		t.Fatalf("erro ao consultar paciente criado: %v", err)
	}

	if nomeSalvo != paciente.Nome {
		t.Fatalf(
			"nome incorreto: esperado %q, recebido %q",
			paciente.Nome,
			nomeSalvo,
		)
	}

	t.Logf(
		"Paciente criado e validado com sucesso: %s",
		paciente.ID,
	)
}

func TestPacienteRepository_DeleteLogico(t *testing.T) {
	db := setupTestDB(t)
	repository := NewPacienteRepository(db)
	ctx := context.Background()

	identificador := uuid.New().String()

	cpf := fmt.Sprintf(
		"%011d",
		time.Now().UnixNano()%100000000000,
	)

	paciente := &models.Paciente{
		Nome: "Paciente Exclusao Logica " + identificador[:8],
		CPF:  cpf,
		DataNascimento: time.Date(
			2000, 1, 15, 0, 0, 0, 0, time.UTC,
		),
		Sexo: models.SexoMasculino,
	}

	err := repository.Create(ctx, paciente)
	if err != nil {
		t.Fatalf("erro ao criar paciente: %v", err)
	}

	t.Cleanup(func() {
		_, err := db.Exec(
			context.Background(),
			"DELETE FROM paciente WHERE id_paciente = $1",
			paciente.ID,
		)
		if err != nil {
			t.Errorf("erro ao limpar paciente: %v", err)
		}
	})

	// Executa a exclusão lógica.
	err = repository.Delete(ctx, paciente.ID)
	if err != nil {
		t.Fatalf("erro na exclusão lógica: %v", err)
	}

	// Confirma que o registro continua no banco.
	var excluidoEm *time.Time

	err = db.QueryRow(ctx, `
		SELECT excluido_em
		FROM paciente
		WHERE id_paciente = $1
	`, paciente.ID).Scan(&excluidoEm)

	if err != nil {
		t.Fatalf("paciente não foi preservado: %v", err)
	}

	if excluidoEm == nil {
		t.Fatal("excluido_em não foi preenchido")
	}

	// Confirma que a consulta operacional não
	// retorna mais o paciente excluído.
	_, err = repository.GetByID(ctx, paciente.ID)
	if err == nil {
		t.Fatal("paciente excluído ainda pode ser consultado")
	}

	// Confirma que a listagem também o oculta.
	pacientes, err := repository.List(ctx)
	if err != nil {
		t.Fatalf("erro ao listar pacientes: %v", err)
	}

	for _, p := range pacientes {
		if p.ID == paciente.ID {
			t.Fatal("paciente excluído apareceu na listagem")
		}
	}

	t.Logf(
		"Exclusão lógica validada: paciente %s preservado e oculto",
		paciente.ID,
	)
}

func TestPacienteRepository_DeleteComFilaAtiva(t *testing.T) {
	statusFila := []string{
		"Aguardando",
		"Chamado",
		"Em Atendimento",
	}

	for _, status := range statusFila {
		t.Run(status, func(t *testing.T) {
			db := setupTestDB(t)
			ctx := context.Background()

			pacienteID, profissionalID := criarDadosTriagemTeste(t, db)

			triagemRepository := NewTriagemRepository(db)
			pacienteRepository := NewPacienteRepository(db)

			triagem := &models.Triagem{
				PacienteID:         pacienteID,
				ProfissionalID:     profissionalID,
				ClassificacaoRisco: models.Medio,
			}

			err := triagemRepository.CreateComFila(ctx, triagem)
			if err != nil {
				t.Fatalf("erro ao criar triagem: %v", err)
			}

			// Altera o status para o cenário testado.
			_, err = db.Exec(ctx, `
				UPDATE fila
				SET status = $1
				WHERE id_triagem = $2
			`, status, triagem.ID)

			if err != nil {
				t.Fatalf("erro ao configurar fila: %v", err)
			}

			// Tenta excluir um paciente com fila ativa.
			err = pacienteRepository.Delete(ctx, pacienteID)

			if !errors.Is(err, ErrPacienteFilaAtiva) {
				t.Fatalf(
					"esperado ErrPacienteFilaAtiva, recebido: %v",
					err,
				)
			}

			// Confirma que o paciente não foi excluído.
			var excluidoEm *time.Time

			err = db.QueryRow(ctx, `
				SELECT excluido_em
				FROM paciente
				WHERE id_paciente = $1
			`, pacienteID).Scan(&excluidoEm)

			if err != nil {
				t.Fatalf("erro ao consultar paciente: %v", err)
			}

			if excluidoEm != nil {
				t.Fatal("paciente com fila ativa foi excluído")
			}

			t.Logf(
				"Exclusão bloqueada corretamente: status %s",
				status,
			)
		})
	}
}

func TestPacienteRepository_DeleteComAtendimentoAtivo(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	pacienteID, profissionalID := criarDadosTriagemTeste(t, db)

	pacienteRepository := NewPacienteRepository(db)
	triagemRepository := NewTriagemRepository(db)

	// Cria uma triagem válida.
	triagem := &models.Triagem{
		PacienteID:         pacienteID,
		ProfissionalID:     profissionalID,
		ClassificacaoRisco: models.Medio,
	}

	err := triagemRepository.CreateComFila(ctx, triagem)
	if err != nil {
		t.Fatalf("erro ao criar triagem: %v", err)
	}

	// Finaliza a participação do paciente na fila,
	// isolando a regra de atendimento clínico ativo.

	// Percorre as transições permitidas pelo trigger
	// até que o paciente não esteja mais na fila ativa.
	for _, status := range []string{
		"Chamado",
		"Em Atendimento",
		"Finalizado",
	} {
		resultado, err := db.Exec(ctx, `
			UPDATE fila
			SET status = $1
			WHERE id_triagem = $2
		`, status, triagem.ID)

		if err != nil {
			t.Fatalf(
				"erro ao alterar fila para %s: %v",
				status,
				err,
			)
		}

		if resultado.RowsAffected() != 1 {
			t.Fatalf(
				"esperada uma fila atualizada para %s, recebidas %d",
				status,
				resultado.RowsAffected(),
			)
		}
	}

	// Obtém a especialidade e o setor do profissional.
	var especialidadeID uuid.UUID
	var setorID uuid.UUID

	err = db.QueryRow(ctx, `
		SELECT id_especialidade, id_setor
		FROM profissional
		WHERE id_profissional = $1
	`, profissionalID).Scan(&especialidadeID, &setorID)

	if err != nil {
		t.Fatalf("erro ao consultar profissional: %v", err)
	}

	// Cria um atendimento clínico em andamento.
	var atendimentoID uuid.UUID

	err = db.QueryRow(ctx, `
		INSERT INTO atendimento (
			id_paciente,
			id_profissional,
			id_triagem,
			id_especialidade,
			id_setor,
			data_inicio,
			data_fim
		)
		VALUES ($1, $2, $3, $4, $5, $6, NULL)
		RETURNING id_atendimento
	`,
		pacienteID,
		profissionalID,
		triagem.ID,
		especialidadeID,
		setorID,
		time.Now(),
	).Scan(&atendimentoID)

	if err != nil {
		t.Fatalf("erro ao criar atendimento: %v", err)
	}

	// Remove o atendimento antes da limpeza dos
	// demais registros realizada pelo helper.
	t.Cleanup(func() {
		_, err := db.Exec(
			context.Background(),
			"DELETE FROM atendimento WHERE id_atendimento = $1",
			atendimentoID,
		)
		if err != nil {
			t.Errorf("erro ao limpar atendimento: %v", err)
		}
	})

	// Tenta excluir o paciente com atendimento aberto.
	err = pacienteRepository.Delete(ctx, pacienteID)

	if !errors.Is(err, ErrPacienteAtendimentoAtivo) {
		t.Fatalf(
			"esperado ErrPacienteAtendimentoAtivo, recebido: %v",
			err,
		)
	}

	// Confirma que o paciente permanece ativo.
	var excluidoEm *time.Time

	err = db.QueryRow(ctx, `
		SELECT excluido_em
		FROM paciente
		WHERE id_paciente = $1
	`, pacienteID).Scan(&excluidoEm)

	if err != nil {
		t.Fatalf("erro ao consultar paciente: %v", err)
	}

	if excluidoEm != nil {
		t.Fatal("paciente com atendimento ativo foi excluído")
	}

	t.Logf(
		"Exclusão bloqueada corretamente: atendimento %s em andamento",
		atendimentoID,
	)
}
