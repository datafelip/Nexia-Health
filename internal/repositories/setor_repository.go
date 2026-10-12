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

var ErrSetorNaoEncontrado = errors.New("setor não encontrado")
var ErrSetorDuplicado = errors.New("já existe um setor com esse nome")
var ErrSetorEmUso = errors.New("setor em uso e não pode ser excluído")

type SetorRepository struct {
	db *pgxpool.Pool
}

func NewSetorRepository(db *pgxpool.Pool) *SetorRepository {
	return &SetorRepository{
		db: db,
	}
}

func setorPgErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

func (r *SetorRepository) Create(ctx context.Context, setor *models.Setor) error {
	query := `
		INSERT INTO setor (
			nome,
			descricao
		)
		VALUES ($1, $2)
		RETURNING id_setor, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		setor.Nome,
		setor.Descricao,
	).Scan(
		&setor.ID,
		&setor.CreatedAt,
		&setor.UpdatedAt,
	)

	if setorPgErrorCode(err) == "23505" { // unique_violation
		return ErrSetorDuplicado
	}

	return err
}

func (r *SetorRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Setor, error) {
	query := `
		SELECT
			id_setor,
			nome,
			descricao,
			created_at,
			updated_at
		FROM setor
		WHERE id_setor = $1
	`

	var setor models.Setor

	err := r.db.QueryRow(ctx, query, id).Scan(
		&setor.ID,
		&setor.Nome,
		&setor.Descricao,
		&setor.CreatedAt,
		&setor.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrSetorNaoEncontrado
	}

	if err != nil {
		return nil, err
	}

	return &setor, nil
}

// List retorna os setores ordenados por nome.
// Se nome != "", filtra por nome contendo o texto (sem diferenciar maiúsculas/minúsculas).
func (r *SetorRepository) List(ctx context.Context, nome string) ([]models.Setor, error) {
	query := `
		SELECT
			id_setor,
			nome,
			descricao,
			created_at,
			updated_at
		FROM setor
		WHERE ($1 = '' OR nome ILIKE '%' || $1 || '%')
		ORDER BY nome ASC
	`

	rows, err := r.db.Query(ctx, query, nome)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	setores := make([]models.Setor, 0)

	for rows.Next() {
		var setor models.Setor

		err := rows.Scan(
			&setor.ID,
			&setor.Nome,
			&setor.Descricao,
			&setor.CreatedAt,
			&setor.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		setores = append(setores, setor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return setores, nil
}

func (r *SetorRepository) Update(ctx context.Context, setor *models.Setor) error {
	query := `
		UPDATE setor
		SET
			nome = $1,
			descricao = $2,
			updated_at = NOW()
		WHERE id_setor = $3
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		setor.Nome,
		setor.Descricao,
		setor.ID,
	).Scan(&setor.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrSetorNaoEncontrado
	}

	if setorPgErrorCode(err) == "23505" { // unique_violation
		return ErrSetorDuplicado
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *SetorRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM setor
		WHERE id_setor = $1
	`

	result, err := r.db.Exec(ctx, query, id)

	if setorPgErrorCode(err) == "23503" { // foreign_key_violation
		return ErrSetorEmUso
	}

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrSetorNaoEncontrado
	}

	return nil
}
