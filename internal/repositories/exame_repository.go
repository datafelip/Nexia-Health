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

var (
	ErrExameNaoEncontrado = errors.New("exame não encontrado")
	ErrExameDuplicado     = errors.New("exame já cadastrado")
	ErrExameEmUso         = errors.New("exame vinculado a outros registros")
)

type ExameRepository struct {
	db *pgxpool.Pool
}

func NewExameRepository(db *pgxpool.Pool) *ExameRepository {
	return &ExameRepository{db: db}
}

func (r *ExameRepository) Create(
	ctx context.Context,
	exame *models.Exame,
) error {
	query := `
		INSERT INTO exame (nome, descricao, preparo)
		VALUES ($1, $2, $3)
		RETURNING id_exame, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		exame.Nome,
		exame.Descricao,
		exame.Preparo,
	).Scan(
		&exame.ID,
		&exame.CreatedAt,
		&exame.UpdatedAt,
	)

	if err != nil {
		return mapExameError(err)
	}

	return nil
}

func (r *ExameRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Exame, error) {
	query := `
		SELECT
			id_exame,
			nome,
			descricao,
			preparo,
			created_at,
			updated_at
		FROM exame
		WHERE id_exame = $1
	`

	var exame models.Exame

	err := r.db.QueryRow(ctx, query, id).Scan(
		&exame.ID,
		&exame.Nome,
		&exame.Descricao,
		&exame.Preparo,
		&exame.CreatedAt,
		&exame.UpdatedAt,
	)

	if err != nil {
		return nil, mapExameError(err)
	}

	return &exame, nil
}

func (r *ExameRepository) List(
	ctx context.Context,
	nome string,
) ([]models.Exame, error) {
	query := `
		SELECT
			id_exame,
			nome,
			descricao,
			preparo,
			created_at,
			updated_at
		FROM exame
		WHERE ($1 = '' OR nome ILIKE '%' || $1 || '%')
		ORDER BY nome
	`

	rows, err := r.db.Query(ctx, query, nome)
	if err != nil {
		return nil, mapExameError(err)
	}
	defer rows.Close()

	exames := make([]models.Exame, 0)

	for rows.Next() {
		var exame models.Exame

		err := rows.Scan(
			&exame.ID,
			&exame.Nome,
			&exame.Descricao,
			&exame.Preparo,
			&exame.CreatedAt,
			&exame.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		exames = append(exames, exame)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return exames, nil
}

func (r *ExameRepository) Update(
	ctx context.Context,
	exame *models.Exame,
) error {
	query := `
		UPDATE exame
		SET
			nome = $2,
			descricao = $3,
			preparo = $4,
			updated_at = NOW()
		WHERE id_exame = $1
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		exame.ID,
		exame.Nome,
		exame.Descricao,
		exame.Preparo,
	).Scan(
		&exame.CreatedAt,
		&exame.UpdatedAt,
	)

	if err != nil {
		return mapExameError(err)
	}

	return nil
}

func (r *ExameRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	result, err := r.db.Exec(
		ctx,
		`DELETE FROM exame WHERE id_exame = $1`,
		id,
	)

	if err != nil {
		return mapExameError(err)
	}

	if result.RowsAffected() == 0 {
		return ErrExameNaoEncontrado
	}

	return nil
}

func mapExameError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrExameNaoEncontrado
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ErrExameDuplicado
		case "23503":
			return ErrExameEmUso
		}
	}

	return err
}
