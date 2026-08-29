package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrTriagemNotFound = errors.New("triagem não encontrada")

type TriagemRepository struct {
	db *pgxpool.Pool
}

func NewTriagemRepository(db *pgxpool.Pool) *TriagemRepository {
	return &TriagemRepository{
		db: db,
	}
}

func (r *TriagemRepository) Create(ctx context.Context, triagem *models.Triagem) error {
	query := `
		INSERT INTO triagem (
			id_paciente,
			id_profissional,
			classificacao_risco,
			pressao_arterial,
			temperatura,
			frequencia_cardiaca,
			observacoes
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			id_triagem,
			data_triagem,
			created_at,
			updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		triagem.PacienteID,
		triagem.ProfissionalID,
		triagem.ClassificacaoRisco,
		triagem.PressaoArterial,
		triagem.Temperatura,
		triagem.FrequenciaCardiaca,
		triagem.Observacoes,
	).Scan(
		&triagem.ID,
		&triagem.DataTriagem,
		&triagem.CreatedAt,
		&triagem.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("erro ao criar triagem: %w", err)
	}

	return nil
}

func (r *TriagemRepository) List(ctx context.Context) ([]models.Triagem, error) {
	query := `
		SELECT
			id_triagem,
			id_paciente,
			id_profissional,
			classificacao_risco,
			pressao_arterial,
			temperatura,
			frequencia_cardiaca,
			observacoes,
			data_triagem,
			created_at,
			updated_at
		FROM triagem
		ORDER BY data_triagem DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar triagens: %w", err)
	}
	defer rows.Close()

	var triagens []models.Triagem

	for rows.Next() {
		var triagem models.Triagem

		err := rows.Scan(
			&triagem.ID,
			&triagem.PacienteID,
			&triagem.ProfissionalID,
			&triagem.ClassificacaoRisco,
			&triagem.PressaoArterial,
			&triagem.Temperatura,
			&triagem.FrequenciaCardiaca,
			&triagem.Observacoes,
			&triagem.DataTriagem,
			&triagem.CreatedAt,
			&triagem.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao ler triagem: %w", err)
		}

		triagens = append(triagens, triagem)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao percorrer triagens: %w", err)
	}

	return triagens, nil
}

func (r *TriagemRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Triagem, error) {
	query := `
		SELECT
			id_triagem,
			id_paciente,
			id_profissional,
			classificacao_risco,
			pressao_arterial,
			temperatura,
			frequencia_cardiaca,
			observacoes,
			data_triagem,
			created_at,
			updated_at
		FROM triagem
		WHERE id_triagem = $1
	`

	var triagem models.Triagem

	err := r.db.QueryRow(ctx, query, id).Scan(
		&triagem.ID,
		&triagem.PacienteID,
		&triagem.ProfissionalID,
		&triagem.ClassificacaoRisco,
		&triagem.PressaoArterial,
		&triagem.Temperatura,
		&triagem.FrequenciaCardiaca,
		&triagem.Observacoes,
		&triagem.DataTriagem,
		&triagem.CreatedAt,
		&triagem.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTriagemNotFound
		}

		return nil, fmt.Errorf("erro ao buscar triagem: %w", err)
	}

	return &triagem, nil
}

func (r *TriagemRepository) Update(ctx context.Context, triagem *models.Triagem) (*models.Triagem, error) {
	query := `
		UPDATE triagem
		SET
			id_paciente = $1,
			id_profissional = $2,
			classificacao_risco = $3,
			pressao_arterial = $4,
			temperatura = $5,
			frequencia_cardiaca = $6,
			observacoes = $7,
			updated_at = NOW()
		WHERE id_triagem = $8
		RETURNING
			id_triagem,
			id_paciente,
			id_profissional,
			classificacao_risco,
			pressao_arterial,
			temperatura,
			frequencia_cardiaca,
			observacoes,
			data_triagem,
			created_at,
			updated_at
	`

	var atualizada models.Triagem

	err := r.db.QueryRow(
		ctx,
		query,
		triagem.PacienteID,
		triagem.ProfissionalID,
		triagem.ClassificacaoRisco,
		triagem.PressaoArterial,
		triagem.Temperatura,
		triagem.FrequenciaCardiaca,
		triagem.Observacoes,
		triagem.ID,
	).Scan(
		&atualizada.ID,
		&atualizada.PacienteID,
		&atualizada.ProfissionalID,
		&atualizada.ClassificacaoRisco,
		&atualizada.PressaoArterial,
		&atualizada.Temperatura,
		&atualizada.FrequenciaCardiaca,
		&atualizada.Observacoes,
		&atualizada.DataTriagem,
		&atualizada.CreatedAt,
		&atualizada.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTriagemNotFound
		}

		return nil, fmt.Errorf("erro ao atualizar triagem: %w", err)
	}

	return &atualizada, nil
}

func (r *TriagemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM triagem
		WHERE id_triagem = $1
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir triagem: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrTriagemNotFound
	}

	return nil
}
