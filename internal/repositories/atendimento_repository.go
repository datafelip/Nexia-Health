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