package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrFilaNaoEncontrada = errors.New("fila não encontrada")

	ErrFilaNaoCancelavel = errors.New(
		"somente pacientes aguardando ou chamados podem ter a fila cancelada",
	)

	ErrFilaVazia = errors.New(
		"não há pacientes aguardando na fila",
	)

	ErrFilaComAtendimento = errors.New(
		"não é permitido cancelar uma fila vinculada a um atendimento",
	)

	ErrFilaOperacaoDiretaNaoPermitida = errors.New(
		"operação direta na fila não permitida; utilize os fluxos integrados",
	)
)

type FilaRepository struct {
	db *pgxpool.Pool
}

func NewFilaRepository(db *pgxpool.Pool) *FilaRepository {
	return &FilaRepository{db: db}
}

// Create permanece disponível por compatibilidade,
// mas impede inserções fora do fluxo integrado de triagem.
func (r *FilaRepository) Create(
	ctx context.Context,
	fila *models.Fila,
) error {
	return ErrFilaOperacaoDiretaNaoPermitida
}

func (r *FilaRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Fila, error) {
	query := `
		SELECT
			id_fila,
			id_paciente,
			id_triagem,
			posicao,
			status,
			data_entrada,
			created_at,
			updated_at
		FROM fila
		WHERE id_fila = $1
	`

	var fila models.Fila

	err := r.db.QueryRow(ctx, query, id).Scan(
		&fila.ID,
		&fila.PacienteID,
		&fila.TriagemID,
		&fila.Posicao,
		&fila.Status,
		&fila.DataEntrada,
		&fila.CreatedAt,
		&fila.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrFilaNaoEncontrada
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar fila: %w", err)
	}

	return &fila, nil
}

func (r *FilaRepository) List(
	ctx context.Context,
) ([]models.Fila, error) {
	query := `
		SELECT
			id_fila,
			id_paciente,
			id_triagem,
			posicao,
			status,
			data_entrada,
			created_at,
			updated_at
		FROM fila
		ORDER BY
			CASE
				WHEN status = 'Aguardando' THEN 0
				ELSE 1
			END,
			posicao ASC,
			data_entrada ASC,
			id_fila ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar fila: %w", err)
	}
	defer rows.Close()

	filas := make([]models.Fila, 0)

	for rows.Next() {
		var fila models.Fila

		err := rows.Scan(
			&fila.ID,
			&fila.PacienteID,
			&fila.TriagemID,
			&fila.Posicao,
			&fila.Status,
			&fila.DataEntrada,
			&fila.CreatedAt,
			&fila.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"erro ao ler registro da fila: %w", err,
			)
		}

		filas = append(filas, fila)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"erro ao percorrer registros da fila: %w", err,
		)
	}

	return filas, nil
}

// Update permanece disponível por compatibilidade,
// mas impede mudanças diretas de status e posição.
func (r *FilaRepository) Update(
	ctx context.Context,
	fila *models.Fila,
) error {
	return ErrFilaOperacaoDiretaNaoPermitida
}

// Delete impede exclusão física do histórico da fila.
func (r *FilaRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	return ErrFilaOperacaoDiretaNaoPermitida
}

// ChamarProximo seleciona o paciente de maior prioridade,
// preservando a ordem de chegada em caso de empate.
func (r *FilaRepository) ChamarProximo(
	ctx context.Context,
) (*models.Fila, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"erro ao iniciar transação: %w", err,
		)
	}
	defer tx.Rollback(ctx)

	// Serializa as operações integradas da fila.
	_, err = tx.Exec(ctx, `
		SELECT pg_advisory_xact_lock(874321)
	`)
	if err != nil {
		return nil, fmt.Errorf(
			"erro ao bloquear fila: %w", err,
		)
	}

	var fila models.Fila

	err = tx.QueryRow(ctx, `
		SELECT
			f.id_fila,
			f.id_paciente,
			f.id_triagem,
			f.posicao,
			f.status,
			f.data_entrada,
			f.created_at,
			f.updated_at
		FROM fila f
		INNER JOIN triagem t
			ON t.id_triagem = f.id_triagem
		WHERE f.status = 'Aguardando'
		  AND t.cancelada_em IS NULL
		ORDER BY
			CASE t.classificacao_risco
				WHEN 'Vermelho' THEN 1
				WHEN 'Amarelo' THEN 2
				WHEN 'Verde' THEN 3
				ELSE 4
			END,
			f.data_entrada ASC,
			f.id_fila ASC
		LIMIT 1
		FOR UPDATE OF f, t
	`).Scan(
		&fila.ID,
		&fila.PacienteID,
		&fila.TriagemID,
		&fila.Posicao,
		&fila.Status,
		&fila.DataEntrada,
		&fila.CreatedAt,
		&fila.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrFilaVazia
	}
	if err != nil {
		return nil, fmt.Errorf(
			"erro ao selecionar próximo paciente: %w", err,
		)
	}

	err = tx.QueryRow(ctx, `
		UPDATE fila
		SET
			status = 'Chamado',
			updated_at = NOW()
		WHERE id_fila = $1
		  AND status = 'Aguardando'
		RETURNING status, updated_at
	`, fila.ID).Scan(
		&fila.Status,
		&fila.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"erro ao chamar paciente: %w", err,
		)
	}

	if err := reordenarFilaAguardando(ctx, tx); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf(
			"erro ao confirmar chamada: %w", err,
		)
	}

	return &fila, nil
}

// CancelarComReordenacao cancela a entrada na fila
// sem apagar a triagem clínica.
//
// Esta operação NÃO preenche triagem.cancelada_em.
// O cancelamento clínico deve utilizar TriagemService.Cancelar.
func (r *FilaRepository) CancelarComReordenacao(
	ctx context.Context,
	id uuid.UUID,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf(
			"erro ao iniciar transação: %w", err,
		)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		SELECT pg_advisory_xact_lock(874321)
	`)
	if err != nil {
		return fmt.Errorf(
			"erro ao bloquear fila: %w", err,
		)
	}

	var status string
	var triagemID uuid.UUID

	err = tx.QueryRow(ctx, `
		SELECT
			f.status::text,
			f.id_triagem
		FROM fila f
		INNER JOIN triagem t
			ON t.id_triagem = f.id_triagem
		WHERE f.id_fila = $1
		FOR UPDATE OF f, t
	`, id).Scan(
		&status,
		&triagemID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrFilaNaoEncontrada
	}
	if err != nil {
		return fmt.Errorf(
			"erro ao consultar fila: %w", err,
		)
	}

	if status != "Aguardando" && status != "Chamado" {
		return ErrFilaNaoCancelavel
	}

	// Verifica se já existe atendimento associado.
	var atendimentoExistente bool

	err = tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM atendimento
			WHERE id_triagem = $1
		)
	`, triagemID).Scan(&atendimentoExistente)

	if err != nil {
		return fmt.Errorf(
			"erro ao verificar atendimentos: %w", err,
		)
	}

	if atendimentoExistente {
		return ErrFilaComAtendimento
	}

	result, err := tx.Exec(ctx, `
		UPDATE fila
		SET
			status = 'Cancelado',
			updated_at = NOW()
		WHERE id_fila = $1
		  AND status IN ('Aguardando', 'Chamado')
	`, id)

	if err != nil {
		return fmt.Errorf(
			"erro ao cancelar fila: %w", err,
		)
	}

	if result.RowsAffected() != 1 {
		return ErrFilaNaoCancelavel
	}

	if err := reordenarFilaAguardando(ctx, tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf(
			"erro ao confirmar cancelamento: %w", err,
		)
	}

	return nil
}

// reordenarFilaAguardando recalcula as posições
// somente dos pacientes que continuam aguardando.
//
// Deve ser chamado dentro de uma transação que já
// tenha adquirido pg_advisory_xact_lock(874321).
func reordenarFilaAguardando(
	ctx context.Context,
	tx pgx.Tx,
) error {
	_, err := tx.Exec(ctx, `
		WITH ordenacao AS (
			SELECT
				f.id_fila,
				ROW_NUMBER() OVER (
					ORDER BY
						CASE t.classificacao_risco
							WHEN 'Vermelho' THEN 1
							WHEN 'Amarelo' THEN 2
							WHEN 'Verde' THEN 3
							ELSE 4
						END,
						f.data_entrada ASC,
						f.id_fila ASC
				)::integer AS nova_posicao
			FROM fila f
			INNER JOIN triagem t
				ON t.id_triagem = f.id_triagem
			WHERE f.status = 'Aguardando'
			  AND t.cancelada_em IS NULL
		)
		UPDATE fila f
		SET
			posicao = o.nova_posicao,
			updated_at = NOW()
		FROM ordenacao o
		WHERE f.id_fila = o.id_fila
		  AND f.posicao IS DISTINCT FROM o.nova_posicao
	`)

	if err != nil {
		return fmt.Errorf(
			"erro ao reorganizar fila: %w", err,
		)
	}

	return nil
}
