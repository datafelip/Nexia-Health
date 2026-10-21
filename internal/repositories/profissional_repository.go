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
	ErrProfissionalNaoEncontrado = errors.New("profissional não encontrado")
	ErrProfissionalDuplicado     = errors.New("CPF ou número de registro já cadastrado")
	ErrProfissionalReferencia    = errors.New("usuário, especialidade ou setor inválido")
)

type ProfissionalRepository struct {
	db *pgxpool.Pool
}

func NewProfissionalRepository(db *pgxpool.Pool) *ProfissionalRepository {
	return &ProfissionalRepository{db: db}
}

func (r *ProfissionalRepository) Create(
	ctx context.Context,
	p *models.Profissional,
) error {
	query := `
		INSERT INTO profissional (
			id_usuario,
			id_especialidade,
			id_setor,
			cpf,
			conselho_profissional,
			numero_registro,
			telefone,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id_profissional, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		p.IDUsuario,
		p.IDEspecialidade,
		p.IDSetor,
		p.CPF,
		p.ConselhoProfissional,
		p.NumeroRegistro,
		p.Telefone,
		p.Status,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return mapProfissionalError(err)
	}

	return nil
}

func (r *ProfissionalRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Profissional, error) {
	query := `
		SELECT
			id_profissional,
			id_usuario,
			id_especialidade,
			id_setor,
			btrim(cpf),
			conselho_profissional,
			numero_registro,
			telefone,
			status,
			created_at,
			updated_at
		FROM profissional
		WHERE id_profissional = $1
	`

	var p models.Profissional

	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.IDUsuario,
		&p.IDEspecialidade,
		&p.IDSetor,
		&p.CPF,
		&p.ConselhoProfissional,
		&p.NumeroRegistro,
		&p.Telefone,
		&p.Status,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		return nil, mapProfissionalError(err)
	}

	return &p, nil
}

func (r *ProfissionalRepository) List(
	ctx context.Context,
	cpf string,
) ([]models.Profissional, error) {
	query := `
		SELECT
			id_profissional,
			id_usuario,
			id_especialidade,
			id_setor,
			btrim(cpf),
			conselho_profissional,
			numero_registro,
			telefone,
			status,
			created_at,
			updated_at
		FROM profissional
		WHERE ($1 = '' OR btrim(cpf) ILIKE '%' || $1 || '%')
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, cpf)
	if err != nil {
		return nil, mapProfissionalError(err)
	}
	defer rows.Close()

	profissionais := make([]models.Profissional, 0)

	for rows.Next() {
		var p models.Profissional

		err := rows.Scan(
			&p.ID,
			&p.IDUsuario,
			&p.IDEspecialidade,
			&p.IDSetor,
			&p.CPF,
			&p.ConselhoProfissional,
			&p.NumeroRegistro,
			&p.Telefone,
			&p.Status,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		profissionais = append(profissionais, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return profissionais, nil
}

func (r *ProfissionalRepository) Update(
	ctx context.Context,
	p *models.Profissional,
) error {
	query := `
		UPDATE profissional
		SET
			id_usuario = $2,
			id_especialidade = $3,
			id_setor = $4,
			cpf = $5,
			conselho_profissional = $6,
			numero_registro = $7,
			telefone = $8,
			status = $9,
			updated_at = NOW()
		WHERE id_profissional = $1
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		p.ID,
		p.IDUsuario,
		p.IDEspecialidade,
		p.IDSetor,
		p.CPF,
		p.ConselhoProfissional,
		p.NumeroRegistro,
		p.Telefone,
		p.Status,
	).Scan(&p.UpdatedAt)

	if err != nil {
		return mapProfissionalError(err)
	}

	return nil
}

func (r *ProfissionalRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	result, err := r.db.Exec(
		ctx,
		`DELETE FROM profissional WHERE id_profissional = $1`,
		id,
	)
	if err != nil {
		return mapProfissionalError(err)
	}

	if result.RowsAffected() == 0 {
		return ErrProfissionalNaoEncontrado
	}

	return nil
}

func mapProfissionalError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrProfissionalNaoEncontrado
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ErrProfissionalDuplicado
		case "23503":
			return ErrProfissionalReferencia
		}
	}

	return err
}
