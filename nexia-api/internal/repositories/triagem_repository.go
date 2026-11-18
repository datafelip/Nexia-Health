package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrTriagemNotFound = errors.New("triagem não encontrada")

var ErrTriagemNaoAguardando = errors.New(
	"somente pacientes aguardando podem ser reclassificados",
)

var ErrTriagemJaCancelada = errors.New(
	"esta triagem já foi cancelada",
)

var ErrTriagemCancelamentoNaoPermitido = errors.New(
	"somente triagens com paciente aguardando podem ser canceladas",
)

var ErrTriagemComAtendimento = errors.New(
	"não é permitido cancelar uma triagem vinculada a um atendimento",
)

var ErrMotivoCancelamentoObrigatorio = errors.New(
	"o motivo do cancelamento é obrigatório",
)

var ErrExclusaoTriagemNaoPermitida = errors.New(
	"a exclusão física de triagens não é permitida",
)

type TriagemRepository struct {
	db *pgxpool.Pool
}

func NewTriagemRepository(db *pgxpool.Pool) *TriagemRepository {
	return &TriagemRepository{
		db: db,
	}
}

// Consulta compartilhada para retornar todos os campos da triagem.
const selectTriagem = `
	SELECT
		id_triagem,
		id_paciente,
		id_profissional,
		classificacao_risco,
		pressao_arterial,
		temperatura,
		frequencia_cardiaca,
		observacoes,
		data_triagem,
		created_at,
		updated_at,
		cancelada_em,
		motivo_cancelamento,
		cancelada_por
	FROM triagem
`

// scanTriagem centraliza a leitura dos campos.
func scanTriagem(
	row pgx.Row,
	triagem *models.Triagem,
) error {
	return row.Scan(
		&triagem.ID,
		&triagem.PacienteID,
		&triagem.ProfissionalID,
		&triagem.ClassificacaoRisco,
		&triagem.PressaoArterial,
		&triagem.Temperatura,
		&triagem.FrequenciaCardiaca,
		&triagem.Observacoes,
		&triagem.DataTriagem,
		&triagem.CreatedAt,
		&triagem.UpdatedAt,
		&triagem.CanceladaEm,
		&triagem.MotivoCancelamento,
		&triagem.CanceladaPor,
	)
}

// Create é mantido para compatibilidade.
// O fluxo HTTP de criação utiliza CreateComFila.
func (r *TriagemRepository) Create(
	ctx context.Context,
	triagem *models.Triagem,
) error {
	query := `
		INSERT INTO triagem (
			id_paciente,
			id_profissional,
			classificacao_risco,
			pressao_arterial,
			temperatura,
			frequencia_cardiaca,
			observacoes
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			id_triagem,
			data_triagem,
			created_at,
			updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		triagem.PacienteID,
		triagem.ProfissionalID,
		triagem.ClassificacaoRisco,
		triagem.PressaoArterial,
		triagem.Temperatura,
		triagem.FrequenciaCardiaca,
		triagem.Observacoes,
	).Scan(
		&triagem.ID,
		&triagem.DataTriagem,
		&triagem.CreatedAt,
		&triagem.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("erro ao criar triagem: %w", err)
	}

	return nil
}

// List retorna todas as triagens, incluindo as canceladas.
// Isso preserva a consulta ao histórico clínico.
func (r *TriagemRepository) List(
	ctx context.Context,
) ([]models.Triagem, error) {
	query := selectTriagem + `
		ORDER BY data_triagem DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar triagens: %w", err)
	}
	defer rows.Close()

	triagens := make([]models.Triagem, 0)

	for rows.Next() {
		var triagem models.Triagem

		err := rows.Scan(
			&triagem.ID,
			&triagem.PacienteID,
			&triagem.ProfissionalID,
			&triagem.ClassificacaoRisco,
			&triagem.PressaoArterial,
			&triagem.Temperatura,
			&triagem.FrequenciaCardiaca,
			&triagem.Observacoes,
			&triagem.DataTriagem,
			&triagem.CreatedAt,
			&triagem.UpdatedAt,
			&triagem.CanceladaEm,
			&triagem.MotivoCancelamento,
			&triagem.CanceladaPor,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao ler triagem: %w", err)
		}

		triagens = append(triagens, triagem)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"erro ao percorrer triagens: %w", err,
		)
	}

	return triagens, nil
}

// GetByID consulta uma triagem independentemente do cancelamento.
func (r *TriagemRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Triagem, error) {
	query := selectTriagem + `
		WHERE id_triagem = $1
	`

	var triagem models.Triagem

	err := scanTriagem(
		r.db.QueryRow(ctx, query, id),
		&triagem,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTriagemNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar triagem: %w", err)
	}

	return &triagem, nil
}

// Update preserva o comportamento anterior,
// impedindo alterações em triagens canceladas.
//
// A rota PUT continua bloqueada no routes.go.
func (r *TriagemRepository) Update(
	ctx context.Context,
	triagem *models.Triagem,
) (*models.Triagem, error) {
	query := `
		UPDATE triagem
		SET
			id_paciente = $1,
			id_profissional = $2,
			classificacao_risco = $3,
			pressao_arterial = $4,
			temperatura = $5,
			frequencia_cardiaca = $6,
			observacoes = $7,
			updated_at = NOW()
		WHERE id_triagem = $8
		  AND cancelada_em IS NULL
		RETURNING
			id_triagem,
			id_paciente,
			id_profissional,
			classificacao_risco,
			pressao_arterial,
			temperatura,
			frequencia_cardiaca,
			observacoes,
			data_triagem,
			created_at,
			updated_at,
			cancelada_em,
			motivo_cancelamento,
			cancelada_por
	`

	var atualizada models.Triagem

	err := scanTriagem(
		r.db.QueryRow(
			ctx,
			query,
			triagem.PacienteID,
			triagem.ProfissionalID,
			triagem.ClassificacaoRisco,
			triagem.PressaoArterial,
			triagem.Temperatura,
			triagem.FrequenciaCardiaca,
			triagem.Observacoes,
			triagem.ID,
		),
		&atualizada,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTriagemNotFound
	}

	if err != nil {
		return nil, fmt.Errorf(
			"erro ao atualizar triagem: %w", err,
		)
	}

	return &atualizada, nil
}

// Delete permanece com a assinatura original,
// mas nunca executa exclusão física.
func (r *TriagemRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	return ErrExclusaoTriagemNaoPermitida
}

// CreateComFila cria a triagem e sua entrada na fila
// dentro da mesma transação.
func (r *TriagemRepository) CreateComFila(
	ctx context.Context,
	triagem *models.Triagem,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf(
			"erro ao iniciar transação: %w", err,
		)
	}
	defer tx.Rollback(ctx)

	// Bloqueio compartilhado entre as operações integradas.
	_, err = tx.Exec(ctx, `
		SELECT pg_advisory_xact_lock(874321)
	`)

	if err != nil {
		return fmt.Errorf(
			"erro ao bloquear fila: %w", err,
		)
	}

	// Impede duas entradas ativas do mesmo paciente.
	var jaExiste bool

	err = tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM fila
			WHERE id_paciente = $1
			  AND status IN (
				'Aguardando',
				'Chamado',
				'Em Atendimento'
			  )
		)
	`, triagem.PacienteID).Scan(&jaExiste)

	if err != nil {
		return fmt.Errorf(
			"erro ao verificar fila: %w", err,
		)
	}

	if jaExiste {
		return fmt.Errorf(
			"paciente já possui entrada ativa na fila",
		)
	}

	// Cria a triagem.
	err = tx.QueryRow(ctx, `
		INSERT INTO triagem (
			id_paciente,
			id_profissional,
			classificacao_risco,
			pressao_arterial,
			temperatura,
			frequencia_cardiaca,
			observacoes
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			id_triagem,
			data_triagem,
			created_at,
			updated_at
	`,
		triagem.PacienteID,
		triagem.ProfissionalID,
		triagem.ClassificacaoRisco,
		triagem.PressaoArterial,
		triagem.Temperatura,
		triagem.FrequenciaCardiaca,
		triagem.Observacoes,
	).Scan(
		&triagem.ID,
		&triagem.DataTriagem,
		&triagem.CreatedAt,
		&triagem.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf(
			"erro ao salvar triagem: %w", err,
		)
	}

	// Insere o paciente na fila.
	_, err = tx.Exec(ctx, `
		INSERT INTO fila (
			id_paciente,
			id_triagem,
			posicao,
			status
		)
		VALUES ($1, $2, 1, 'Aguardando')
	`, triagem.PacienteID, triagem.ID)

	if err != nil {
		return fmt.Errorf(
			"erro ao inserir na fila: %w", err,
		)
	}

	// Reorganiza a fila por classificação de risco.
	if err := reordenarFilaAguardando(ctx, tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf(
			"erro ao confirmar transação: %w", err,
		)
	}

	return nil
}

// ReclassificarComFila altera a classificação de risco
// e reorganiza a fila na mesma transação.
func (r *TriagemRepository) ReclassificarComFila(
	ctx context.Context,
	id uuid.UUID,
	novaClassificacao models.ClassificacaoRisco,
) (*models.Triagem, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"erro ao iniciar transação: %w", err,
		)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		SELECT pg_advisory_xact_lock(874321)
	`)

	if err != nil {
		return nil, fmt.Errorf(
			"erro ao bloquear fila: %w", err,
		)
	}

	// Diferencia triagem inexistente de triagem
	// que não pode mais ser reclassificada.
	var existe bool

	err = tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM triagem
			WHERE id_triagem = $1
		)
	`, id).Scan(&existe)

	if err != nil {
		return nil, fmt.Errorf(
			"erro ao verificar triagem: %w", err,
		)
	}

	if !existe {
		return nil, ErrTriagemNotFound
	}

	// Somente pacientes aguardando podem ser reclassificados.
	var filaID uuid.UUID

	err = tx.QueryRow(ctx, `
		SELECT f.id_fila
		FROM fila f
		INNER JOIN triagem t
			ON t.id_triagem = f.id_triagem
		WHERE f.id_triagem = $1
		  AND f.status = 'Aguardando'
		  AND t.cancelada_em IS NULL
		FOR UPDATE OF f, t
	`, id).Scan(&filaID)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTriagemNaoAguardando
	}

	if err != nil {
		return nil, fmt.Errorf(
			"erro ao consultar fila: %w", err,
		)
	}

	result, err := tx.Exec(ctx, `
		UPDATE triagem
		SET
			classificacao_risco = $1,
			updated_at = NOW()
		WHERE id_triagem = $2
		  AND cancelada_em IS NULL
	`, novaClassificacao, id)

	if err != nil {
		return nil, fmt.Errorf(
			"erro ao reclassificar triagem: %w", err,
		)
	}

	if result.RowsAffected() != 1 {
		return nil, ErrTriagemNaoAguardando
	}

	if err := reordenarFilaAguardando(ctx, tx); err != nil {
		return nil, err
	}

	// Retorna a triagem atualizada.
	var triagem models.Triagem

	err = scanTriagem(
		tx.QueryRow(
			ctx,
			selectTriagem+` WHERE id_triagem = $1`,
			id,
		),
		&triagem,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"erro ao buscar triagem atualizada: %w", err,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf(
			"erro ao confirmar reclassificação: %w", err,
		)
	}

	return &triagem, nil
}

// CancelarComFila cancela a triagem e a entrada na fila,
// preservando o histórico clínico e a justificativa.
//
// Toda a operação ocorre em uma única transação.
func (r *TriagemRepository) CancelarComFila(
	ctx context.Context,
	id uuid.UUID,
	motivo string,
) (*models.Triagem, error) {
	motivo = strings.TrimSpace(motivo)

	if motivo == "" {
		return nil, ErrMotivoCancelamentoObrigatorio
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"erro ao iniciar transação: %w", err,
		)
	}
	defer tx.Rollback(ctx)

	// Mesmo bloqueio utilizado pela fila e pelos atendimentos.
	_, err = tx.Exec(ctx, `
		SELECT pg_advisory_xact_lock(874321)
	`)

	if err != nil {
		return nil, fmt.Errorf(
			"erro ao bloquear fila: %w", err,
		)
	}

	// Bloqueia a triagem durante o cancelamento.
	var canceladaEm *time.Time

	err = tx.QueryRow(ctx, `
		SELECT cancelada_em
		FROM triagem
		WHERE id_triagem = $1
		FOR UPDATE
	`, id).Scan(&canceladaEm)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTriagemNotFound
	}

	if err != nil {
		return nil, fmt.Errorf(
			"erro ao consultar triagem: %w", err,
		)
	}

	if canceladaEm != nil {
		return nil, ErrTriagemJaCancelada
	}

	// Impede cancelamento quando existe atendimento vinculado.
	var possuiAtendimento bool

	err = tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM atendimento
			WHERE id_triagem = $1
		)
	`, id).Scan(&possuiAtendimento)

	if err != nil {
		return nil, fmt.Errorf(
			"erro ao verificar atendimento vinculado: %w",
			err,
		)
	}

	if possuiAtendimento {
		return nil, ErrTriagemComAtendimento
	}

	// Apenas pacientes aguardando podem ter a triagem cancelada.
	var filaID uuid.UUID

	err = tx.QueryRow(ctx, `
		SELECT id_fila
		FROM fila
		WHERE id_triagem = $1
		  AND status = 'Aguardando'
		FOR UPDATE
	`, id).Scan(&filaID)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTriagemCancelamentoNaoPermitido
	}

	if err != nil {
		return nil, fmt.Errorf(
			"erro ao consultar fila da triagem: %w",
			err,
		)
	}

	// Registra o cancelamento sem excluir dados clínicos.
	// cancelada_por permanece NULL até haver autenticação.
	resultadoTriagem, err := tx.Exec(ctx, `
		UPDATE triagem
		SET
			cancelada_em = NOW(),
			motivo_cancelamento = $1,
			updated_at = NOW()
		WHERE id_triagem = $2
		  AND cancelada_em IS NULL
	`, motivo, id)

	if err != nil {
		return nil, fmt.Errorf(
			"erro ao registrar cancelamento: %w", err,
		)
	}

	if resultadoTriagem.RowsAffected() != 1 {
		return nil, ErrTriagemJaCancelada
	}

	// Cancela a entrada correspondente na fila.
	resultadoFila, err := tx.Exec(ctx, `
		UPDATE fila
		SET
			status = 'Cancelado',
			updated_at = NOW()
		WHERE id_fila = $1
		  AND status = 'Aguardando'
	`, filaID)

	if err != nil {
		return nil, fmt.Errorf(
			"erro ao cancelar entrada na fila: %w", err,
		)
	}

	if resultadoFila.RowsAffected() != 1 {
		return nil, ErrTriagemCancelamentoNaoPermitido
	}

	// Reorganiza os pacientes que permanecem aguardando.
	if err := reordenarFilaAguardando(ctx, tx); err != nil {
		return nil, err
	}

	// Recupera os dados atualizados para a resposta HTTP.
	var triagem models.Triagem

	err = scanTriagem(
		tx.QueryRow(
			ctx,
			selectTriagem+` WHERE id_triagem = $1`,
			id,
		),
		&triagem,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"erro ao buscar triagem cancelada: %w", err,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf(
			"erro ao confirmar cancelamento: %w", err,
		)
	}

	return &triagem, nil
}
