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
	ErrSolicitacaoNaoEncontrada = errors.New("solicitação de exame não encontrada")
	ErrSolicitacaoReferencia    = errors.New("atendimento ou exame inexistente")
	ErrSolicitacaoDuplicada     = errors.New("solicitação de exame duplicada")
	ErrSolicitacaoEmUso         = errors.New("solicitação de exame vinculada a outros registros")
)

type SolicitacaoExameRepository struct {
	db *pgxpool.Pool
}

func NewSolicitacaoExameRepository(db *pgxpool.Pool) *SolicitacaoExameRepository {
	return &SolicitacaoExameRepository{db: db}
}

func (r *SolicitacaoExameRepository) Create(
	ctx context.Context,
	s *models.SolicitacaoExame,
) error {
	query := `
		INSERT INTO solicitacao_exame (
			id_atendimento,
			id_exame,
			status,
			resultado
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id_solicitacao, data_solicitacao, created_at, updated_at
	`

	return mapSolicitacaoError(r.db.QueryRow(
		ctx, query,
		s.IDAtendimento,
		s.IDExame,
		s.Status,
		s.Resultado,
	).Scan(
		&s.ID,
		&s.DataSolicitacao,
		&s.CreatedAt,
		&s.UpdatedAt,
	))
}

func (r *SolicitacaoExameRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.SolicitacaoExame, error) {
	query := `
		SELECT id_solicitacao, id_atendimento, id_exame,
		       data_solicitacao, status, resultado,
		       created_at, updated_at
		FROM solicitacao_exame
		WHERE id_solicitacao = $1
	`

	var s models.SolicitacaoExame

	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID,
		&s.IDAtendimento,
		&s.IDExame,
		&s.DataSolicitacao,
		&s.Status,
		&s.Resultado,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		return nil, mapSolicitacaoError(err)
	}

	return &s, nil
}

func (r *SolicitacaoExameRepository) List(
	ctx context.Context,
	atendimentoID *uuid.UUID,
) ([]models.SolicitacaoExame, error) {
	query := `
		SELECT id_solicitacao, id_atendimento, id_exame,
		       data_solicitacao, status, resultado,
		       created_at, updated_at
		FROM solicitacao_exame
		WHERE ($1::uuid IS NULL OR id_atendimento = $1::uuid)
		ORDER BY data_solicitacao DESC
	`

	rows, err := r.db.Query(ctx, query, atendimentoID)
	if err != nil {
		return nil, mapSolicitacaoError(err)
	}
	defer rows.Close()

	solicitacoes := make([]models.SolicitacaoExame, 0)

	for rows.Next() {
		var s models.SolicitacaoExame

		if err := rows.Scan(
			&s.ID,
			&s.IDAtendimento,
			&s.IDExame,
			&s.DataSolicitacao,
			&s.Status,
			&s.Resultado,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}

		solicitacoes = append(solicitacoes, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return solicitacoes, nil
}

func (r *SolicitacaoExameRepository) Update(
	ctx context.Context,
	s *models.SolicitacaoExame,
) error {
	query := `
		UPDATE solicitacao_exame
		SET id_atendimento = $2,
		    id_exame = $3,
		    status = $4,
		    resultado = $5,
		    updated_at = NOW()
		WHERE id_solicitacao = $1
		RETURNING data_solicitacao, created_at, updated_at
	`

	return mapSolicitacaoError(r.db.QueryRow(
		ctx, query,
		s.ID,
		s.IDAtendimento,
		s.IDExame,
		s.Status,
		s.Resultado,
	).Scan(
		&s.DataSolicitacao,
		&s.CreatedAt,
		&s.UpdatedAt,
	))
}

func (r *SolicitacaoExameRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	result, err := r.db.Exec(
		ctx,
		`DELETE FROM solicitacao_exame WHERE id_solicitacao = $1`,
		id,
	)
	if err != nil {
		return mapSolicitacaoError(err)
	}

	if result.RowsAffected() == 0 {
		return ErrSolicitacaoNaoEncontrada
	}

	return nil
}

func mapSolicitacaoError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrSolicitacaoNaoEncontrada
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ErrSolicitacaoDuplicada
		case "23503":
			if pgErr.ConstraintName == "fk_solicitacao_atendimento" ||
				pgErr.ConstraintName == "fk_solicitacao_exame" {
				return ErrSolicitacaoReferencia
			}
			return ErrSolicitacaoEmUso
		}
	}

	return err
}
