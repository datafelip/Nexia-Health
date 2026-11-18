package repositories

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// criarDadosTriagemTeste cria registros fictícios
// exclusivamente no banco de testes.
//
// Todos os registros são removidos automaticamente
// ao final do teste.
func criarDadosTriagemTeste(
	t *testing.T,
	db *pgxpool.Pool,
) (uuid.UUID, uuid.UUID) {
	t.Helper()

	ctx := context.Background()
	identificador := uuid.New().String()

	var pacienteID uuid.UUID
	var usuarioID uuid.UUID
	var especialidadeID uuid.UUID
	var setorID uuid.UUID
	var profissionalID uuid.UUID

	// A limpeza é registrada antes das inserções,
	// para também funcionar em caso de falha parcial.
	t.Cleanup(func() {
		ctx := context.Background()

		// A ordem respeita as chaves estrangeiras.
		consultas := []struct {
			query string
			id    uuid.UUID
		}{
			{"DELETE FROM fila WHERE id_paciente = $1", pacienteID},
			{"DELETE FROM triagem WHERE id_paciente = $1", pacienteID},
			{"DELETE FROM profissional WHERE id_profissional = $1", profissionalID},
			{"DELETE FROM paciente WHERE id_paciente = $1", pacienteID},
			{"DELETE FROM usuario WHERE id_usuario = $1", usuarioID},
			{"DELETE FROM especialidade WHERE id_especialidade = $1", especialidadeID},
			{"DELETE FROM setor WHERE id_setor = $1", setorID},
		}

		for _, consulta := range consultas {
			if consulta.id == uuid.Nil {
				continue
			}

			if _, err := db.Exec(ctx, consulta.query, consulta.id); err != nil {
				t.Errorf("erro na limpeza dos dados de teste: %v", err)
			}
		}
	})

	// Paciente fictício.
	cpfPaciente := fmt.Sprintf(
		"%011d",
		time.Now().UnixNano()%100000000000,
	)

	err := db.QueryRow(ctx, `
		INSERT INTO paciente (
			nome,
			cpf,
			data_nascimento,
			sexo
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id_paciente
	`,
		"Paciente Teste "+identificador[:8],
		cpfPaciente,
		time.Date(2000, 1, 15, 0, 0, 0, 0, time.UTC),
		"Masculino",
	).Scan(&pacienteID)

	if err != nil {
		t.Fatalf("erro ao criar paciente de teste: %v", err)
	}

	// Usuário fictício associado ao profissional.
	err = db.QueryRow(ctx, `
		INSERT INTO usuario (
			nome,
			email,
			senha_hash,
			perfil
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id_usuario
	`,
		"Enfermeiro Teste",
		"teste-"+identificador+"@example.invalid",
		"hash-ficticio-sem-acesso",
		"Enfermeiro",
	).Scan(&usuarioID)

	if err != nil {
		t.Fatalf("erro ao criar usuário de teste: %v", err)
	}

	// Especialidade fictícia.
	err = db.QueryRow(ctx, `
		INSERT INTO especialidade (nome)
		VALUES ($1)
		RETURNING id_especialidade
	`,
		"Especialidade Teste "+identificador,
	).Scan(&especialidadeID)

	if err != nil {
		t.Fatalf("erro ao criar especialidade: %v", err)
	}

	// Setor fictício.
	err = db.QueryRow(ctx, `
		INSERT INTO setor (nome)
		VALUES ($1)
		RETURNING id_setor
	`,
		"Setor Teste "+identificador,
	).Scan(&setorID)

	if err != nil {
		t.Fatalf("erro ao criar setor: %v", err)
	}

	// Profissional fictício.
	cpfProfissional := fmt.Sprintf(
		"%011d",
		(time.Now().UnixNano()+1)%100000000000,
	)

	err = db.QueryRow(ctx, `
		INSERT INTO profissional (
			id_usuario,
			id_especialidade,
			id_setor,
			cpf,
			conselho_profissional,
			numero_registro
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id_profissional
	`,
		usuarioID,
		especialidadeID,
		setorID,
		cpfProfissional,
		"COREN",
		"TESTE-"+identificador[:8],
	).Scan(&profissionalID)

	if err != nil {
		t.Fatalf("erro ao criar profissional de teste: %v", err)
	}

	return pacienteID, profissionalID
}

func TestTriagemRepository_CreateComFila(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	pacienteID, profissionalID := criarDadosTriagemTeste(t, db)

	repository := NewTriagemRepository(db)

	pressao := "120/80"
	temperatura := 36.5
	frequencia := int16(78)
	observacoes := "Triagem fictícia para teste automatizado"

	triagem := &models.Triagem{
		PacienteID:         pacienteID,
		ProfissionalID:     profissionalID,
		ClassificacaoRisco: models.Medio,
		PressaoArterial:    &pressao,
		Temperatura:        &temperatura,
		FrequenciaCardiaca: &frequencia,
		Observacoes:        &observacoes,
	}

	err := repository.CreateComFila(ctx, triagem)
	if err != nil {
		t.Fatalf("erro ao criar triagem com fila: %v", err)
	}

	if triagem.ID == uuid.Nil {
		t.Fatal("UUID da triagem não foi gerado")
	}

	var status string
	var posicao int

	err = db.QueryRow(ctx, `
		SELECT status::text, posicao
		FROM fila
		WHERE id_triagem = $1
	`, triagem.ID).Scan(&status, &posicao)

	if err != nil {
		t.Fatalf("erro ao consultar fila criada: %v", err)
	}

	if status != "Aguardando" {
		t.Fatalf(
			"status incorreto: esperado Aguardando, recebido %s",
			status,
		)
	}

	if posicao < 1 {
		t.Fatalf("posição inválida na fila: %d", posicao)
	}

	t.Logf(
		"Triagem %s criada com fila na posição %d",
		triagem.ID,
		posicao,
	)
}

func TestTriagemRepository_CancelarComFila(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	pacienteID, profissionalID := criarDadosTriagemTeste(t, db)

	repository := NewTriagemRepository(db)

	observacoes := "Paciente avaliado durante a triagem"

	triagem := &models.Triagem{
		PacienteID:         pacienteID,
		ProfissionalID:     profissionalID,
		ClassificacaoRisco: models.Medio,
		Observacoes:        &observacoes,
	}

	err := repository.CreateComFila(ctx, triagem)
	if err != nil {
		t.Fatalf("erro ao preparar triagem: %v", err)
	}

	motivo := "Cadastro realizado em duplicidade"

	triagemCancelada, err := repository.CancelarComFila(
		ctx,
		triagem.ID,
		motivo,
	)

	if err != nil {
		t.Fatalf("erro ao cancelar triagem: %v", err)
	}

	if triagemCancelada == nil {
		t.Fatal("cancelamento retornou triagem nula")
	}

	if triagemCancelada.CanceladaEm == nil {
		t.Fatal("data do cancelamento não foi registrada")
	}

	if triagemCancelada.MotivoCancelamento == nil {
		t.Fatal("motivo do cancelamento não foi registrado")
	}

	if *triagemCancelada.MotivoCancelamento != motivo {
		t.Fatalf(
			"motivo incorreto: esperado %q, recebido %q",
			motivo,
			*triagemCancelada.MotivoCancelamento,
		)
	}

	// Confirma que o histórico clínico permanece no banco.
	triagemSalva, err := repository.GetByID(ctx, triagem.ID)
	if err != nil {
		t.Fatalf("triagem não encontrada após cancelamento: %v", err)
	}

	if triagemSalva == nil || triagemSalva.ID != triagem.ID {
		t.Fatal("histórico da triagem não foi preservado")
	}

	if triagemSalva.CanceladaEm == nil {
		t.Fatal("cancelamento não persistiu no banco")
	}

	if triagemSalva.Observacoes == nil ||
		*triagemSalva.Observacoes != observacoes {
		t.Fatal("observações clínicas foram alteradas")
	}

	// Confirma que a fila foi cancelada.
	var status string

	err = db.QueryRow(ctx, `
		SELECT status::text
		FROM fila
		WHERE id_triagem = $1
	`, triagem.ID).Scan(&status)

	if err != nil {
		t.Fatalf("erro ao consultar fila cancelada: %v", err)
	}

	if status != "Cancelado" {
		t.Fatalf(
			"status incorreto: esperado Cancelado, recebido %s",
			status,
		)
	}

	// Uma segunda tentativa não deve alterar o histórico.
	_, err = repository.CancelarComFila(
		ctx,
		triagem.ID,
		"Segunda tentativa de cancelamento",
	)

	if !errors.Is(err, ErrTriagemJaCancelada) {
		t.Fatalf(
			"esperado ErrTriagemJaCancelada, recebido: %v",
			err,
		)
	}

	// Verifica se o motivo original foi preservado.
	var motivoPersistido string

	err = db.QueryRow(ctx, `
		SELECT motivo_cancelamento
		FROM triagem
		WHERE id_triagem = $1
	`, triagem.ID).Scan(&motivoPersistido)

	if err != nil {
		t.Fatalf("erro ao consultar motivo persistido: %v", err)
	}

	if motivoPersistido != motivo {
		t.Fatal("segunda tentativa alterou o motivo original")
	}

	t.Logf(
		"Cancelamento validado com histórico preservado: %s",
		triagem.ID,
	)
}

func TestTriagemRepository_CancelarSemMotivo(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	repository := NewTriagemRepository(db)

	_, err := repository.CancelarComFila(
		ctx,
		uuid.New(),
		"   ",
	)

	if !errors.Is(err, ErrMotivoCancelamentoObrigatorio) {
		t.Fatalf(
			"esperado ErrMotivoCancelamentoObrigatorio, recebido: %v",
			err,
		)
	}
}

func TestTriagemRepository_ReordenarFilaAposCancelamento(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	repository := NewTriagemRepository(db)

	// Cada paciente possui dados próprios e independentes.
	pacienteA, profissionalA := criarDadosTriagemTeste(t, db)
	pacienteB, profissionalB := criarDadosTriagemTeste(t, db)
	pacienteC, profissionalC := criarDadosTriagemTeste(t, db)

	triagemA := &models.Triagem{
		PacienteID:         pacienteA,
		ProfissionalID:     profissionalA,
		ClassificacaoRisco: models.Medio,
	}

	triagemB := &models.Triagem{
		PacienteID:         pacienteB,
		ProfissionalID:     profissionalB,
		ClassificacaoRisco: models.Medio,
	}

	triagemC := &models.Triagem{
		PacienteID:         pacienteC,
		ProfissionalID:     profissionalC,
		ClassificacaoRisco: models.Baixo,
	}

	// Cria as triagens na ordem A, B, C.
	for _, triagem := range []*models.Triagem{
		triagemA,
		triagemB,
		triagemC,
	} {
		if err := repository.CreateComFila(ctx, triagem); err != nil {
			t.Fatalf("erro ao criar triagem: %v", err)
		}
	}

	// Consulta a posição de uma triagem na fila.
	consultarFila := func(
		id uuid.UUID,
	) (int, string) {
		t.Helper()

		var posicao int
		var status string

		err := db.QueryRow(ctx, `
			SELECT posicao, status::text
			FROM fila
			WHERE id_triagem = $1
		`, id).Scan(&posicao, &status)

		if err != nil {
			t.Fatalf("erro ao consultar fila: %v", err)
		}

		return posicao, status
	}

	// Confirma a ordem inicial.
	posA, statusA := consultarFila(triagemA.ID)
	posB, statusB := consultarFila(triagemB.ID)
	posC, statusC := consultarFila(triagemC.ID)

	if posA != 1 || posB != 2 || posC != 3 {
		t.Fatalf(
			"ordem inicial incorreta: A=%d, B=%d, C=%d",
			posA, posB, posC,
		)
	}

	if statusA != "Aguardando" ||
		statusB != "Aguardando" ||
		statusC != "Aguardando" {
		t.Fatal("todos os pacientes deveriam estar aguardando")
	}

	// Cancela a triagem intermediária.
	_, err := repository.CancelarComFila(
		ctx,
		triagemB.ID,
		"Cancelamento para testar reorganização",
	)

	if err != nil {
		t.Fatalf("erro ao cancelar triagem B: %v", err)
	}

	// Consulta novamente a fila.
	posA, statusA = consultarFila(triagemA.ID)
	posB, statusB = consultarFila(triagemB.ID)
	posC, statusC = consultarFila(triagemC.ID)

	if posA != 1 || statusA != "Aguardando" {
		t.Fatalf(
			"paciente A incorreto: posição=%d, status=%s",
			posA, statusA,
		)
	}

	if statusB != "Cancelado" {
		t.Fatalf(
			"paciente B deveria estar cancelado: %s",
			statusB,
		)
	}

	if posC != 2 || statusC != "Aguardando" {
		t.Fatalf(
			"paciente C incorreto: posição=%d, status=%s",
			posC, statusC,
		)
	}

	// Confirma que a triagem B permanece no banco.
	triagemCancelada, err := repository.GetByID(
		ctx,
		triagemB.ID,
	)

	if err != nil {
		t.Fatalf(
			"erro ao recuperar histórico da triagem B: %v",
			err,
		)
	}

	if triagemCancelada == nil ||
		triagemCancelada.CanceladaEm == nil {
		t.Fatal("histórico de cancelamento não preservado")
	}

	// Verifica que não existem lacunas na fila ativa.
	var quantidadeAguardando int
	var maiorPosicao int

	err = db.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COALESCE(MAX(posicao), 0)
		FROM fila
		WHERE status = 'Aguardando'
	`).Scan(&quantidadeAguardando, &maiorPosicao)

	if err != nil {
		t.Fatalf("erro ao verificar posições da fila: %v", err)
	}

	if quantidadeAguardando != 2 || maiorPosicao != 2 {
		t.Fatalf(
			"fila inconsistente: quantidade=%d, maior posição=%d",
			quantidadeAguardando,
			maiorPosicao,
		)
	}

	t.Log(
		"Fila reorganizada corretamente após cancelamento: " +
			"A=1, B=Cancelado, C=2",
	)
}

func TestTriagemRepository_CancelarTriagemInexistente(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	repository := NewTriagemRepository(db)

	// UUID aleatório que não corresponde
	// a nenhuma triagem cadastrada.
	triagemID := uuid.New()

	triagem, err := repository.CancelarComFila(
		ctx,
		triagemID,
		"Tentativa de cancelar triagem inexistente",
	)

	if !errors.Is(err, ErrTriagemNotFound) {
		t.Fatalf(
			"esperado ErrTriagemNotFound, recebido: %v",
			err,
		)
	}

	if triagem != nil {
		t.Fatal("triagem inexistente não deveria ser retornada")
	}

	t.Log("Cancelamento de triagem inexistente bloqueado corretamente")
}

func TestTriagemRepository_CancelarPacienteJaChamado(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	triagemRepository := NewTriagemRepository(db)
	filaRepository := NewFilaRepository(db)

	// Cria os dados fictícios do paciente e profissional.
	pacienteID, profissionalID := criarDadosTriagemTeste(t, db)

	triagem := &models.Triagem{
		PacienteID:         pacienteID,
		ProfissionalID:     profissionalID,
		ClassificacaoRisco: models.Medio,
	}

	// 1. Cria a triagem e coloca o paciente na fila.
	err := triagemRepository.CreateComFila(ctx, triagem)
	if err != nil {
		t.Fatalf("erro ao criar triagem: %v", err)
	}

	// 2. Chama o próximo paciente pelo fluxo real.
	filaChamada, err := filaRepository.ChamarProximo(ctx)
	if err != nil {
		t.Fatalf("erro ao chamar paciente: %v", err)
	}

	if filaChamada == nil {
		t.Fatal("nenhuma fila retornada após chamada")
	}

	if filaChamada.TriagemID != triagem.ID {
		t.Fatalf(
			"triagem chamada incorreta: esperado %s, recebido %s",
			triagem.ID,
			filaChamada.TriagemID,
		)
	}

	if string(filaChamada.Status) != "Chamado" {
		t.Fatalf(
			"status esperado Chamado, recebido %s",
			filaChamada.Status,
		)
	}

	// 3. Tenta cancelar a triagem já chamada.
	triagemCancelada, err := triagemRepository.CancelarComFila(
		ctx,
		triagem.ID,
		"Tentativa de cancelamento após chamada",
	)

	if !errors.Is(err, ErrTriagemCancelamentoNaoPermitido) {
		t.Fatalf(
			"esperado ErrTriagemCancelamentoNaoPermitido, recebido: %v",
			err,
		)
	}

	if triagemCancelada != nil {
		t.Fatal("cancelamento bloqueado não deveria retornar triagem")
	}

	// 4. Verifica se o status da fila permaneceu Chamado.
	var status string

	err = db.QueryRow(ctx, `
		SELECT status::text
		FROM fila
		WHERE id_triagem = $1
	`, triagem.ID).Scan(&status)

	if err != nil {
		t.Fatalf("erro ao consultar status da fila: %v", err)
	}

	if status != "Chamado" {
		t.Fatalf(
			"status alterado indevidamente: esperado Chamado, recebido %s",
			status,
		)
	}

	// 5. Confirma que a triagem não foi cancelada.
	triagemSalva, err := triagemRepository.GetByID(
		ctx,
		triagem.ID,
	)

	if err != nil {
		t.Fatalf("erro ao recuperar triagem: %v", err)
	}

	if triagemSalva == nil {
		t.Fatal("triagem não encontrada")
	}

	if triagemSalva.CanceladaEm != nil {
		t.Fatal("triagem chamada foi cancelada indevidamente")
	}

	if triagemSalva.MotivoCancelamento != nil {
		t.Fatal("motivo de cancelamento registrado indevidamente")
	}

	t.Log(
		"Cancelamento bloqueado corretamente: " +
			"paciente já chamado permanece com status Chamado",
	)
}

func TestTriagemRepository_CancelarComAtendimentoVinculado(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	triagemRepository := NewTriagemRepository(db)
	atendimentoRepository := NewAtendimentoRepository(db)

	// 1. Cria paciente e profissional fictícios.
	pacienteID, profissionalID := criarDadosTriagemTeste(t, db)

	// 2. Recupera especialidade e setor do profissional.
	var especialidadeID uuid.UUID
	var setorID uuid.UUID

	err := db.QueryRow(ctx, `
		SELECT id_especialidade, id_setor
		FROM profissional
		WHERE id_profissional = $1
	`, profissionalID).Scan(&especialidadeID, &setorID)

	if err != nil {
		t.Fatalf(
			"erro ao recuperar especialidade e setor: %v",
			err,
		)
	}

	// 3. Cria a triagem com sua entrada na fila.
	triagem := &models.Triagem{
		PacienteID:         pacienteID,
		ProfissionalID:     profissionalID,
		ClassificacaoRisco: models.Medio,
	}

	err = triagemRepository.CreateComFila(ctx, triagem)
	if err != nil {
		t.Fatalf("erro ao criar triagem: %v", err)
	}

	// 4. Prepara o atendimento.
	observacoes := "Atendimento fictício para teste automatizado"

	atendimento := &models.Atendimento{
		PacienteID:      pacienteID,
		ProfissionalID:  profissionalID,
		TriagemID:       triagem.ID,
		EspecialidadeID: especialidadeID,
		SetorID:         setorID,
		DataInicio:      time.Now(),
		Observacoes:     &observacoes,
	}

	// Registra a limpeza do atendimento antes da limpeza
	// dos dados de paciente e profissional.
	t.Cleanup(func() {
		_, err := db.Exec(
			context.Background(),
			"DELETE FROM atendimento WHERE id_triagem = $1",
			triagem.ID,
		)
		if err != nil {
			t.Errorf("erro ao limpar atendimento de teste: %v", err)
		}
	})

	// 5. Cria o atendimento pelo fluxo real.
	err = atendimentoRepository.CreateComFila(ctx, atendimento)
	if err != nil {
		t.Fatalf("erro ao criar atendimento: %v", err)
	}

	if atendimento.ID == uuid.Nil {
		t.Fatal("UUID do atendimento não foi gerado")
	}

	// 6. Confirma o status Em Atendimento.
	var statusAntes string

	err = db.QueryRow(ctx, `
		SELECT status::text
		FROM fila
		WHERE id_triagem = $1
	`, triagem.ID).Scan(&statusAntes)

	if err != nil {
		t.Fatalf("erro ao consultar fila: %v", err)
	}

	if statusAntes != "Em Atendimento" {
		t.Fatalf(
			"esperado Em Atendimento, recebido %s",
			statusAntes,
		)
	}

	// 7. Tenta cancelar a triagem vinculada ao atendimento.
	triagemCancelada, err := triagemRepository.CancelarComFila(
		ctx,
		triagem.ID,
		"Tentativa de cancelamento com atendimento iniciado",
	)

	if !errors.Is(err, ErrTriagemComAtendimento) {
		t.Fatalf(
			"esperado ErrTriagemComAtendimento, recebido: %v",
			err,
		)
	}

	if triagemCancelada != nil {
		t.Fatal("cancelamento recusado não deveria retornar triagem")
	}

	// 8. Confirma que o atendimento continua armazenado.
	var quantidadeAtendimentos int

	err = db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM atendimento
		WHERE id_atendimento = $1
		  AND id_triagem = $2
	`, atendimento.ID, triagem.ID).Scan(&quantidadeAtendimentos)

	if err != nil {
		t.Fatalf("erro ao consultar atendimento: %v", err)
	}

	if quantidadeAtendimentos != 1 {
		t.Fatal("atendimento foi removido indevidamente")
	}

	// 9. Confirma que a fila não foi alterada.
	var statusDepois string

	err = db.QueryRow(ctx, `
		SELECT status::text
		FROM fila
		WHERE id_triagem = $1
	`, triagem.ID).Scan(&statusDepois)

	if err != nil {
		t.Fatalf("erro ao consultar status após cancelamento: %v", err)
	}

	if statusDepois != "Em Atendimento" {
		t.Fatalf(
			"status alterado indevidamente: %s",
			statusDepois,
		)
	}

	// 10. Confirma que a triagem permanece sem cancelamento.
	triagemSalva, err := triagemRepository.GetByID(
		ctx,
		triagem.ID,
	)

	if err != nil {
		t.Fatalf("erro ao recuperar triagem: %v", err)
	}

	if triagemSalva == nil {
		t.Fatal("triagem não encontrada após tentativa de cancelamento")
	}

	if triagemSalva.CanceladaEm != nil {
		t.Fatal("triagem com atendimento foi cancelada indevidamente")
	}

	if triagemSalva.MotivoCancelamento != nil {
		t.Fatal("motivo de cancelamento registrado indevidamente")
	}

	t.Log(
		"Cancelamento bloqueado corretamente: " +
			"atendimento preservado e fila Em Atendimento",
	)
}

func TestTriagemRepository_CancelamentoConcorrente(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	repository := NewTriagemRepository(db)

	pacienteID, profissionalID := criarDadosTriagemTeste(t, db)

	triagem := &models.Triagem{
		PacienteID:         pacienteID,
		ProfissionalID:     profissionalID,
		ClassificacaoRisco: models.Medio,
	}

	if err := repository.CreateComFila(ctx, triagem); err != nil {
		t.Fatalf("erro ao criar triagem: %v", err)
	}

	type resultado struct {
		motivo string
		err    error
	}

	resultados := make(chan resultado, 2)
	iniciar := make(chan struct{})

	motivos := []string{
		"Cancelamento solicitado pelo profissional A",
		"Cancelamento solicitado pelo profissional B",
	}

	for _, motivo := range motivos {
		motivo := motivo

		go func() {
			// As duas goroutines aguardam o mesmo sinal.
			<-iniciar

			_, err := repository.CancelarComFila(
				ctx,
				triagem.ID,
				motivo,
			)

			resultados <- resultado{
				motivo: motivo,
				err:    err,
			}
		}()
	}

	// Libera as duas tentativas de cancelamento.
	close(iniciar)

	sucessos := 0
	jaCanceladas := 0
	motivoVencedor := ""

	for i := 0; i < 2; i++ {
		r := <-resultados

		switch {
		case r.err == nil:
			sucessos++
			motivoVencedor = r.motivo

		case errors.Is(r.err, ErrTriagemJaCancelada):
			jaCanceladas++

		default:
			t.Fatalf(
				"erro inesperado no cancelamento concorrente: %v",
				r.err,
			)
		}
	}

	if sucessos != 1 {
		t.Fatalf(
			"esperado exatamente 1 cancelamento bem-sucedido, recebido %d",
			sucessos,
		)
	}

	if jaCanceladas != 1 {
		t.Fatalf(
			"esperado exatamente 1 ErrTriagemJaCancelada, recebido %d",
			jaCanceladas,
		)
	}

	// Confirma que o motivo vencedor foi preservado.
	triagemSalva, err := repository.GetByID(ctx, triagem.ID)
	if err != nil {
		t.Fatalf("erro ao consultar triagem: %v", err)
	}

	if triagemSalva == nil || triagemSalva.CanceladaEm == nil {
		t.Fatal("cancelamento não foi persistido")
	}

	if triagemSalva.MotivoCancelamento == nil {
		t.Fatal("motivo do cancelamento não foi persistido")
	}

	if *triagemSalva.MotivoCancelamento != motivoVencedor {
		t.Fatalf(
			"motivo original sobrescrito: esperado %q, recebido %q",
			motivoVencedor,
			*triagemSalva.MotivoCancelamento,
		)
	}

	// Confirma que existe somente uma entrada de fila
	// associada à triagem e que ela está cancelada.
	var quantidade int
	var status string

	err = db.QueryRow(ctx, `
		SELECT COUNT(*), MAX(status::text)
		FROM fila
		WHERE id_triagem = $1
	`, triagem.ID).Scan(&quantidade, &status)

	if err != nil {
		t.Fatalf("erro ao consultar fila: %v", err)
	}

	if quantidade != 1 || status != "Cancelado" {
		t.Fatalf(
			"fila inconsistente: quantidade=%d, status=%s",
			quantidade,
			status,
		)
	}

	t.Log(
		"Concorrência validada: uma operação cancelou " +
			"e a outra foi bloqueada sem sobrescrever o histórico",
	)
}
