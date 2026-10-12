package repositories

import (
	"context"
	"errors"
	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrEspecialidadeNaoEncontrada = errors.New("especialidade não encontrada")
var ErrEspecialidadeDuplicada = errors.New("já existe uma especialidade com esse nome")
var ErrEspecialidadeEmUso = errors.New("especialidade em uso e não pode ser excluída")

type EspecialidadeRepository struct {
	db *pgxpool.Pool
}

func NewEspecialidadeRepository(db *pgxpool.Pool) *EspecialidadeRepository {
	return &EspecialidadeRepository{
		db: db,
	}
}

// pgErrorCode devolve o código SQLSTATE do erro do Postgres (ou "" se não for um).
func pgErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

func (r *EspecialidadeRepository) Create(ctx context.Context, especialidade *models.Especialidade) error {
	query := `
		INSERT INTO especialidade (
			nome,
			descricao
		)
		VALUES ($1, $2)
		RETURNING id_especialidade, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		especialidade.Nome,
		especialidade.Descricao,
	).Scan(
		&especialidade.ID,
		&especialidade.CreatedAt,
		&especialidade.UpdatedAt,
	)

	if pgErrorCode(err) == "23505" { // unique_violation
		return ErrEspecialidadeDuplicada
	}

	return err
}

func (r *EspecialidadeRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Especialidade, error) {
	query := `
		SELECT
			id_especialidade,
			nome,
			descricao,
			created_at,
			updated_at
		FROM especialidade
		WHERE id_especialidade = $1
	`

	var especialidade models.Especialidade

	err := r.db.QueryRow(ctx, query, id).Scan(
		&especialidade.ID,
		&especialidade.Nome,
		&especialidade.Descricao,
		&especialidade.CreatedAt,
		&especialidade.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrEspecialidadeNaoEncontrada
	}

	if err != nil {
		return nil, err
	}

	return &especialidade, nil
}

// List retorna as especialidades ordenadas por nome.
// Se nome != "", filtra por nome contendo o texto (sem diferenciar maiúsculas/minúsculas).
func (r *EspecialidadeRepository) List(ctx context.Context, nome string) ([]models.Especialidade, error) {
	query := `
		SELECT
			id_especialidade,
			nome,
			descricao,
			created_at,
			updated_at
		FROM especialidade
		WHERE ($1 = '' OR nome ILIKE '%' || $1 || '%')
		ORDER BY nome ASC
	`

	rows, err := r.db.Query(ctx, query, nome)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	especialidades := make([]models.Especialidade, 0)

	for rows.Next() {
		var especialidade models.Especialidade

		err := rows.Scan(
			&especialidade.ID,
			&especialidade.Nome,
			&especialidade.Descricao,
			&especialidade.CreatedAt,
			&especialidade.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		especialidades = append(especialidades, especialidade)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return especialidades, nil
}

func (r *EspecialidadeRepository) Update(ctx context.Context, especialidade *models.Especialidade) error {
	query := `
		UPDATE especialidade
		SET
			nome = $1,
			descricao = $2,
			updated_at = NOW()
		WHERE id_especialidade = $3
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		especialidade.Nome,
		especialidade.Descricao,
		especialidade.ID,
	).Scan(&especialidade.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrEspecialidadeNaoEncontrada
	}

	if pgErrorCode(err) == "23505" { // unique_violation
		return ErrEspecialidadeDuplicada
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *EspecialidadeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM especialidade
		WHERE id_especialidade = $1
	`

	result, err := r.db.Exec(ctx, query, id)

	if pgErrorCode(err) == "23503" { // foreign_key_violation
		return ErrEspecialidadeEmUso
	}

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrEspecialidadeNaoEncontrada
	}

	return nil
}
