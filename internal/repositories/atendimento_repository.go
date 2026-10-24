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

var ErrAtendimentoNaoEncontrado = errors.New("atendimento não encontrado")

type AtendimentoRepository struct {
	db *pgxpool.Pool
}

func NewAtendimentoRepository(db *pgxpool.Pool) *AtendimentoRepository {
	return &AtendimentoRepository{
		db: db,
	}
}

// AtendimentoFilter representa os filtros opcionais aceitos por List.
// Campos nil/zero são ignorados na query.
type AtendimentoFilter struct {
	PacienteID *uuid.UUID
	DataInicio *time.Time // limite inferior (inclusive) sobre atendimento.data_inicio
	DataFim    *time.Time // limite superior (inclusive) sobre atendimento.data_inicio
}

func (r *AtendimentoRepository) Create(ctx context.Context, atendimento *models.Atendimento) error {
	query := `
		INSERT INTO atendimento (
			id_paciente,
			id_profissional,
			id_triagem,
			id_agendamento,
			id_especialidade,
			id_setor,
			data_inicio,
			data_fim,
			diagnostico,
			observacoes
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
		RETURNING id_atendimento, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		atendimento.PacienteID,
		atendimento.ProfissionalID,
		atendimento.TriagemID,
		atendimento.AgendamentoID,
		atendimento.EspecialidadeID,
		atendimento.SetorID,
		atendimento.DataInicio,
		atendimento.DataFim,
		atendimento.Diagnostico,
		atendimento.Observacoes,
	).Scan(
		&atendimento.ID,
		&atendimento.CreatedAt,
		&atendimento.UpdatedAt,
	)

	return err
}

func (r *AtendimentoRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Atendimento, error) {
	query := `
		SELECT
			id_atendimento,
			id_paciente,
			id_profissional,
			id_triagem,
			id_agendamento,
			id_especialidade,
			id_setor,
			data_inicio,
			data_fim,
			diagnostico,
			observacoes,
			created_at,
			updated_at
		FROM atendimento
		WHERE id_atendimento = $1
	`

	var atendimento models.Atendimento

	err := r.db.QueryRow(ctx, query, id).Scan(
		&atendimento.ID,
		&atendimento.PacienteID,
		&atendimento.ProfissionalID,
		&atendimento.TriagemID,
		&atendimento.AgendamentoID,
		&atendimento.EspecialidadeID,
		&atendimento.SetorID,
		&atendimento.DataInicio,
		&atendimento.DataFim,
		&atendimento.Diagnostico,
		&atendimento.Observacoes,
		&atendimento.CreatedAt,
		&atendimento.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrAtendimentoNaoEncontrado
	}

	if err != nil {
		return nil, err
	}

	return &atendimento, nil
}

// List retorna atendimentos aplicando os filtros informados em filter.
// Um filter zerado (todos os campos nil) equivale a listar tudo.
func (r *AtendimentoRepository) List(ctx context.Context, filter AtendimentoFilter) ([]models.Atendimento, error) {
	baseQuery := `
		SELECT
			id_atendimento,
			id_paciente,
			id_profissional,
			id_triagem,
			id_agendamento,
			id_especialidade,
			id_setor,
			data_inicio,
			data_fim,
			diagnostico,
			observacoes,
			created_at,
			updated_at
		FROM atendimento
	`

	var conditions []string
	var args []interface{}
	argPos := 1

	if filter.PacienteID != nil {
		conditions = append(conditions, fmt.Sprintf("id_paciente = $%d", argPos))
		args = append(args, *filter.PacienteID)
		argPos++
	}

	if filter.DataInicio != nil {
		conditions = append(conditions, fmt.Sprintf("data_inicio >= $%d", argPos))
		args = append(args, *filter.DataInicio)
		argPos++
	}

	if filter.DataFim != nil {
		conditions = append(conditions, fmt.Sprintf("data_inicio <= $%d", argPos))
		args = append(args, *filter.DataFim)
		argPos++
	}

	query := baseQuery
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY data_inicio DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	atendimentos := make([]models.Atendimento, 0)

	for rows.Next() {
		var atendimento models.Atendimento

		err := rows.Scan(
			&atendimento.ID,
			&atendimento.PacienteID,
			&atendimento.ProfissionalID,
			&atendimento.TriagemID,
			&atendimento.AgendamentoID,
			&atendimento.EspecialidadeID,
			&atendimento.SetorID,
			&atendimento.DataInicio,
			&atendimento.DataFim,
			&atendimento.Diagnostico,
			&atendimento.Observacoes,
			&atendimento.CreatedAt,
			&atendimento.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		atendimentos = append(atendimentos, atendimento)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return atendimentos, nil
}

func (r *AtendimentoRepository) Update(ctx context.Context, atendimento *models.Atendimento) error {
	query := `
		UPDATE atendimento
		SET
			data_fim = $1,
			diagnostico = $2,
			observacoes = $3,
			updated_at = NOW()
		WHERE id_atendimento = $4
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		atendimento.DataFim,
		atendimento.Diagnostico,
		atendimento.Observacoes,
		atendimento.ID,
	).Scan(&atendimento.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrAtendimentoNaoEncontrado
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *AtendimentoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM atendimento
		WHERE id_atendimento = $1
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrAtendimentoNaoEncontrado
	}

	return nil
}

func (r *AtendimentoRepository) CreateComFila(
	ctx context.Context,
	atendimento *models.Atendimento,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	// Mesmo bloqueio utilizado na criação de triagem com fila.
	_, err = tx.Exec(ctx,
		"SELECT pg_advisory_xact_lock(874321)",
	)
	if err != nil {
		return fmt.Errorf("erro ao bloquear fila: %w", err)
	}

	// Verifica se o paciente possui entrada ativa
	// correspondente à triagem informada.
	var filaID uuid.UUID
	var status models.Status

	err = tx.QueryRow(ctx, `
		SELECT id_fila, status
		FROM fila
		WHERE id_paciente = $1
		  AND id_triagem = $2
		  AND status IN ('Aguardando', 'Chamado')
		FOR UPDATE
	`, atendimento.PacienteID, atendimento.TriagemID).
		Scan(&filaID, &status)

	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf(
			"paciente não possui entrada aguardando ou chamada para esta triagem",
		)
	}
	if err != nil {
		return fmt.Errorf("erro ao consultar fila: %w", err)
	}

	// Impede reutilizar uma triagem em outro atendimento.
	var atendimentoExistente bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM atendimento
			WHERE id_triagem = $1
		)
	`, atendimento.TriagemID).Scan(&atendimentoExistente)
	if err != nil {
		return fmt.Errorf("erro ao verificar triagem: %w", err)
	}
	if atendimentoExistente {
		return fmt.Errorf("esta triagem já possui atendimento")
	}

	// Impede outro atendimento aberto para o mesmo paciente.
	var atendimentoAberto bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM atendimento
			WHERE id_paciente = $1
			  AND data_fim IS NULL
		)
	`, atendimento.PacienteID).Scan(&atendimentoAberto)
	if err != nil {
		return fmt.Errorf("erro ao verificar atendimentos: %w", err)
	}
	if atendimentoAberto {
		return fmt.Errorf("paciente já possui atendimento aberto")
	}

	// Cria o atendimento.
	err = tx.QueryRow(ctx, `
		INSERT INTO atendimento (
			id_paciente,
			id_profissional,
			id_triagem,
			id_agendamento,
			id_especialidade,
			id_setor,
			data_inicio,
			data_fim,
			diagnostico,
			observacoes
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id_atendimento, created_at, updated_at
	`,
		atendimento.PacienteID,
		atendimento.ProfissionalID,
		atendimento.TriagemID,
		atendimento.AgendamentoID,
		atendimento.EspecialidadeID,
		atendimento.SetorID,
		atendimento.DataInicio,
		atendimento.DataFim,
		atendimento.Diagnostico,
		atendimento.Observacoes,
	).Scan(
		&atendimento.ID,
		&atendimento.CreatedAt,
		&atendimento.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar atendimento: %w", err)
	}

	// Atualiza o status da fila.
	_, err = tx.Exec(ctx, `
		UPDATE fila
		SET status = 'Em Atendimento',
		    updated_at = NOW()
		WHERE id_fila = $1
	`, filaID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar status da fila: %w", err)
	}

	// Recalcula somente as posições de quem aguarda.
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
						f.data_entrada,
						f.id_fila
				)::integer AS nova_posicao
			FROM fila f
			JOIN triagem t ON t.id_triagem = f.id_triagem
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
		return fmt.Errorf("erro ao confirmar atendimento: %w", err)
	}

	return nil
}

func (r *AtendimentoRepository) FinalizarComFila(
	ctx context.Context,
	atendimento *models.Atendimento,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	// Utiliza o mesmo bloqueio dos outros fluxos integrados.
	_, err = tx.Exec(ctx,
		"SELECT pg_advisory_xact_lock(874321)",
	)
	if err != nil {
		return fmt.Errorf("erro ao bloquear fila: %w", err)
	}

	// Bloqueia e verifica o atendimento atual.
	var triagemID uuid.UUID
	var dataInicio time.Time
	var dataFimAtual *time.Time

	err = tx.QueryRow(ctx, `
		SELECT id_triagem, data_inicio, data_fim
		FROM atendimento
		WHERE id_atendimento = $1
		FOR UPDATE
	`, atendimento.ID).Scan(
		&triagemID,
		&dataInicio,
		&dataFimAtual,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrAtendimentoNaoEncontrado
	}
	if err != nil {
		return fmt.Errorf("erro ao consultar atendimento: %w", err)
	}

	if dataFimAtual != nil {
		return fmt.Errorf("atendimento já está finalizado")
	}

	if atendimento.DataFim == nil {
		return fmt.Errorf("data_fim é obrigatória para finalizar")
	}

	if atendimento.DataFim.Before(dataInicio) {
		return fmt.Errorf("data_fim não pode ser anterior a data_inicio")
	}

	// Atualiza o atendimento.
	_, err = tx.Exec(ctx, `
		UPDATE atendimento
		SET data_fim = $1,
		    diagnostico = $2,
		    observacoes = $3,
		    updated_at = NOW()
		WHERE id_atendimento = $4
	`,
		atendimento.DataFim,
		atendimento.Diagnostico,
		atendimento.Observacoes,
		atendimento.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao finalizar atendimento: %w", err)
	}

	// Finaliza a entrada correspondente na fila.
	result, err := tx.Exec(ctx, `
		UPDATE fila
		SET status = 'Finalizado',
		    updated_at = NOW()
		WHERE id_triagem = $1
		  AND status = 'Em Atendimento'
	`, triagemID)
	if err != nil {
		return fmt.Errorf("erro ao finalizar fila: %w", err)
	}

	if result.RowsAffected() != 1 {
		return fmt.Errorf(
			"esperada uma entrada Em Atendimento para esta triagem; encontradas %d",
			result.RowsAffected(),
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("erro ao confirmar finalização: %w", err)
	}

	return nil
}
