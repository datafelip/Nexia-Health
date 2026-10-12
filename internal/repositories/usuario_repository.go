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

var ErrUsuarioNaoEncontrado = errors.New("usuário não encontrado")
var ErrUsuarioDuplicado = errors.New("já existe um usuário com esse e-mail")

type UsuarioRepository struct {
	db *pgxpool.Pool
}

func NewUsuarioRepository(db *pgxpool.Pool) *UsuarioRepository {
	return &UsuarioRepository{
		db: db,
	}
}

func usuarioPgErrorCode(err error) string {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		return pgErr.Code
	}

	return ""
}

func (r *UsuarioRepository) Create(ctx context.Context, usuario *models.Usuario) error {
	query := `
		INSERT INTO usuario (
			nome,
			email,
			senha_hash,
			perfil,
			status,
			ultimo_login
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id_usuario, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		usuario.Nome,
		usuario.Email,
		usuario.SenhaHash,
		usuario.Perfil,
		usuario.Status,
		usuario.UltimoLogin,
	).Scan(
		&usuario.ID,
		&usuario.CreatedAt,
		&usuario.UpdatedAt,
	)

	if usuarioPgErrorCode(err) == "23505" {
		return ErrUsuarioDuplicado
	}

	return err
}

func (r *UsuarioRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Usuario, error) {
	query := `
		SELECT
			id_usuario,
			nome,
			email,
			senha_hash,
			perfil,
			status,
			ultimo_login,
			created_at,
			updated_at
		FROM usuario
		WHERE id_usuario = $1
	`

	var usuario models.Usuario

	err := r.db.QueryRow(ctx, query, id).Scan(
		&usuario.ID,
		&usuario.Nome,
		&usuario.Email,
		&usuario.SenhaHash,
		&usuario.Perfil,
		&usuario.Status,
		&usuario.UltimoLogin,
		&usuario.CreatedAt,
		&usuario.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUsuarioNaoEncontrado
	}

	if err != nil {
		return nil, err
	}

	return &usuario, nil
}

func (r *UsuarioRepository) List(
	ctx context.Context,
	nome string,
) ([]models.Usuario, error) {
	query := `
		SELECT
			id_usuario,
			nome,
			email,
			senha_hash,
			perfil,
			status,
			ultimo_login,
			created_at,
			updated_at
		FROM usuario
		WHERE ($1 = '' OR nome ILIKE '%' || $1 || '%')
		ORDER BY nome ASC
	`

	rows, err := r.db.Query(ctx, query, nome)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usuarios := make([]models.Usuario, 0)

	for rows.Next() {
		var usuario models.Usuario

		err := rows.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Email,
			&usuario.SenhaHash,
			&usuario.Perfil,
			&usuario.Status,
			&usuario.UltimoLogin,
			&usuario.CreatedAt,
			&usuario.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		usuarios = append(usuarios, usuario)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return usuarios, nil
}

func (r *UsuarioRepository) Update(
	ctx context.Context,
	usuario *models.Usuario,
) error {
	query := `
		UPDATE usuario
		SET
			nome = $1,
			email = $2,
			senha_hash = $3,
			perfil = $4,
			status = $5,
			ultimo_login = $6,
			updated_at = NOW()
		WHERE id_usuario = $7
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		usuario.Nome,
		usuario.Email,
		usuario.SenhaHash,
		usuario.Perfil,
		usuario.Status,
		usuario.UltimoLogin,
		usuario.ID,
	).Scan(&usuario.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrUsuarioNaoEncontrado
	}

	if usuarioPgErrorCode(err) == "23505" {
		return ErrUsuarioDuplicado
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *UsuarioRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	query := `
		DELETE FROM usuario
		WHERE id_usuario = $1
	`

	result, err := r.db.Exec(ctx, query, id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrUsuarioNaoEncontrado
	}

	return nil
}
