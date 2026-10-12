package repositories

import (
	"context"
	"errors"
	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrFilaNaoEncontrada = errors.New("fila não encontrada")

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