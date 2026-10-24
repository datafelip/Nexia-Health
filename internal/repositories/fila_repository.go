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

var ErrFilaNaoEncontrada = errors.New("fila não encontrada")

var ErrFilaNaoCancelavel = errors.New(
	"somente pacientes aguardando ou chamados podem ter a fila cancelada",
)

type FilaRepository struct {
	db *pgxpool.Pool
}

func NewFilaRepository(db *pgxpool.Pool) *FilaRepository {
	return &FilaRepository{
		db: db,
	}
}

func (r *FilaRepository) Create(ctx context.Context, fila *models.Fila) error {
	query := `
		INSERT INTO fila (
			id_paciente,
			id_triagem,
			posicao,
			status,
			data_entrada
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id_fila, created_at, updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		fila.PacienteID,
		fila.TriagemID,
		fila.Posicao,
		fila.Status,
		fila.DataEntrada,
	).Scan(
		&fila.ID,
		&fila.CreatedAt,
		&fila.UpdatedAt,
	)
}

func (r *FilaRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Fila, error) {
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
		return nil, err
	}

	return &fila, nil
}

func (r *FilaRepository) List(ctx context.Context) ([]models.Fila, error) {
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
		ORDER BY posicao ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		filas = append(filas, fila)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return filas, nil
}

func (r *FilaRepository) Update(ctx context.Context, fila *models.Fila) error {
	query := `
		UPDATE fila
		SET
			posicao = $1,
			status = $2,
			updated_at = NOW()
		WHERE id_fila = $3
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		fila.Posicao,
		fila.Status,
		fila.ID,
	).Scan(&fila.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrFilaNaoEncontrada
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *FilaRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM fila
		WHERE id_fila = $1
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrFilaNaoEncontrada
	}

	return nil
}

var ErrFilaVazia = errors.New(
	"não há pacientes aguardando na fila",
)

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

	// Mesmo bloqueio utilizado nas integrações anteriores.
	// Impede que duas chamadas deste fluxo selecionem
	// simultaneamente o mesmo próximo paciente.
	_, err = tx.Exec(ctx,
		"SELECT pg_advisory_xact_lock(874321)",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"erro ao bloquear fila: %w", err,
		)
	}

	var fila models.Fila

	// Seleciona o paciente de maior prioridade.
	// Em caso de empate, respeita a ordem de chegada.
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
		JOIN triagem t
			ON t.id_triagem = f.id_triagem
		WHERE f.status = 'Aguardando'
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
		FOR UPDATE OF f
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
			"erro ao selecionar paciente: %w", err,
		)
	}

	// Altera somente o paciente selecionado.
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

	// Reorganiza apenas os pacientes que continuam
	// aguardando, sem interromper os atendimentos.
	_, err = tx.Exec(ctx, `
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
			JOIN triagem t
				ON t.id_triagem = f.id_triagem
			WHERE f.status = 'Aguardando'
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
		return nil, fmt.Errorf(
			"erro ao reorganizar fila: %w", err,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf(
			"erro ao confirmar chamada: %w", err,
		)
	}

	return &fila, nil
}

func (r *FilaRepository) CancelarComReordenacao(
	ctx context.Context,
	id uuid.UUID,
) error {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	// Mesmo bloqueio utilizado pelas demais operações da fila.
	_, err = tx.Exec(ctx, `
		SELECT pg_advisory_xact_lock(874321)
	`)
	if err != nil {
		return fmt.Errorf("erro ao bloquear fila: %w", err)
	}

	// Consulta e bloqueia o registro da fila.
	var status string

	err = tx.QueryRow(ctx, `
		SELECT status::text
		FROM fila
		WHERE id_fila = $1
		FOR UPDATE
	`, id).Scan(&status)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrFilaNaoEncontrada
	}

	if err != nil {
		return fmt.Errorf("erro ao consultar fila: %w", err)
	}

	// Somente esses estados permitem cancelamento.
	if status != "Aguardando" && status != "Chamado" {
		return ErrFilaNaoCancelavel
	}

	// Atualiza o status sem excluir o histórico.
	_, err = tx.Exec(ctx, `
		UPDATE fila
		SET status = 'Cancelado',
		    updated_at = NOW()
		WHERE id_fila = $1
	`, id)

	if err != nil {
		return fmt.Errorf("erro ao cancelar fila: %w", err)
	}

	// Reorganiza os pacientes que continuam aguardando.
	_, err = tx.Exec(ctx, `
		WITH ordenacao AS (
			SELECT
				f.id_fila,
				ROW_NUMBER() OVER (
					ORDER BY
						CASE t.classificacao_risco
							WHEN 'Vermelho' THEN 1
							WHEN 'Amarelo' THEN 2
							WHEN 'Verde' THEN 3
						END,
						f.data_entrada ASC,
						f.id_fila ASC
				)::integer AS nova_posicao
			FROM fila f
			JOIN triagem t
			  ON t.id_triagem = f.id_triagem
			WHERE f.status = 'Aguardando'
		)
		UPDATE fila f
		SET posicao = o.nova_posicao,
		    updated_at = NOW()
		FROM ordenacao o
		WHERE f.id_fila = o.id_fila
		  AND f.posicao IS DISTINCT FROM o.nova_posicao
	`)

	if err != nil {
		return fmt.Errorf("erro ao reorganizar fila: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("erro ao confirmar cancelamento: %w", err)
	}

	return nil
}
