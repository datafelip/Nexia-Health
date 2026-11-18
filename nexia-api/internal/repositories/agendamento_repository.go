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

var ErrAgendamentoNaoEncontrado = errors.New("agendamento não encontrado")

type AgendamentoRepository struct {
	db *pgxpool.Pool
}

func NewAgendamentoRepository(db *pgxpool.Pool) *AgendamentoRepository {
	return &AgendamentoRepository{
		db: db,
	}
}

// AgendamentoFilter representa os filtros opcionais aceitos por List.
// Campos nil são ignorados na query.
type AgendamentoFilter struct {
	PacienteID *uuid.UUID
	DataInicio *time.Time // limite inferior (inclusive) sobre data_agendada
	DataFim    *time.Time // limite superior (inclusive) sobre data_agendada
	Status     *models.StatusAgendamento
}

func (r *AgendamentoRepository) Create(ctx context.Context, agendamento *models.Agendamento) error {
	query := `
		INSERT INTO agendamento (
			id_paciente,
			tipo,
			data_agendada,
			horario,
			status
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id_agendamento, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		agendamento.PacienteID,
		agendamento.Tipo,
		agendamento.DataHora.Format("2006-01-02"), // coluna data_agendada (DATE)
		agendamento.DataHora.Format("15:04:05"),    // coluna horario (TIME)
		agendamento.Status,
	).Scan(
		&agendamento.ID,
		&agendamento.CreatedAt,
		&agendamento.UpdatedAt,
	)

	return err
}

func (r *AgendamentoRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Agendamento, error) {
	query := `
		SELECT
			id_agendamento,
			id_paciente,
			tipo,
			(data_agendada + horario) AS data_hora,
			status,
			created_at,
			updated_at
		FROM agendamento
		WHERE id_agendamento = $1
	`

	var agendamento models.Agendamento

	err := r.db.QueryRow(ctx, query, id).Scan(
		&agendamento.ID,
		&agendamento.PacienteID,
		&agendamento.Tipo,
		&agendamento.DataHora,
		&agendamento.Status,
		&agendamento.CreatedAt,
		&agendamento.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrAgendamentoNaoEncontrado
	}

	if err != nil {
		return nil, err
	}

	return &agendamento, nil
}

// List retorna agendamentos aplicando os filtros informados em filter.
// Um filter zerado (todos os campos nil) equivale a listar tudo.
func (r *AgendamentoRepository) List(ctx context.Context, filter AgendamentoFilter) ([]models.Agendamento, error) {
	baseQuery := `
		SELECT
			id_agendamento,
			id_paciente,
			tipo,
			(data_agendada + horario) AS data_hora,
			status,
			created_at,
			updated_at
		FROM agendamento
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
		conditions = append(conditions, fmt.Sprintf("data_agendada >= $%d", argPos))
		args = append(args, filter.DataInicio.Format("2006-01-02"))
		argPos++
	}

	if filter.DataFim != nil {
		conditions = append(conditions, fmt.Sprintf("data_agendada <= $%d", argPos))
		args = append(args, filter.DataFim.Format("2006-01-02"))
		argPos++
	}

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argPos))
		args = append(args, *filter.Status)
		argPos++
	}

	query := baseQuery
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY data_agendada ASC, horario ASC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	agendamentos := make([]models.Agendamento, 0)

	for rows.Next() {
		var agendamento models.Agendamento

		err := rows.Scan(
			&agendamento.ID,
			&agendamento.PacienteID,
			&agendamento.Tipo,
			&agendamento.DataHora,
			&agendamento.Status,
			&agendamento.CreatedAt,
			&agendamento.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		agendamentos = append(agendamentos, agendamento)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return agendamentos, nil
}

func (r *AgendamentoRepository) Update(ctx context.Context, agendamento *models.Agendamento) error {
	query := `
		UPDATE agendamento
		SET
			tipo = $1,
			data_agendada = $2,
			horario = $3,
			status = $4,
			updated_at = NOW()
		WHERE id_agendamento = $5
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		agendamento.Tipo,
		agendamento.DataHora.Format("2006-01-02"),
		agendamento.DataHora.Format("15:04:05"),
		agendamento.Status,
		agendamento.ID,
	).Scan(&agendamento.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrAgendamentoNaoEncontrado
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *AgendamentoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM agendamento
		WHERE id_agendamento = $1
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrAgendamentoNaoEncontrado
	}

	return nil
}