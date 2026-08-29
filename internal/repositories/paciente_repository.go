package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/felipelima/nexia-health/api/internal/models"
)

var ErrPacienteNotFound = errors.New("Paciente não encontrado")

type PacienteRepository struct {
	db *pgxpool.Pool
}

func NewPacienteRepository(db *pgxpool.Pool) *PacienteRepository {
	return &PacienteRepository{
		db: db,
	}
}

func (r *PacienteRepository) Create(ctx context.Context, paciente *models.Paciente) error {
	query := `
		INSERT INTO paciente (
			id_paciente,
			nome,
			cpf,
			data_nascimento,
			sexo,
			telefone,
			endereco,
			cartao_sus
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		RETURNING data_cadastro, created_at, updated_at
	`

	if paciente.ID == uuid.Nil {
		paciente.ID = uuid.New()
	}

	err := r.db.QueryRow(
		ctx,
		query,
		paciente.ID,
		paciente.Nome,
		paciente.CPF,
		paciente.DataNascimento,
		paciente.Sexo,
		paciente.Telefone,
		paciente.Endereco,
		paciente.CartaoSUS,
	).Scan(
		&paciente.DataCadastro,
		&paciente.CreatedAt,
		&paciente.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("erro ao criar paciente: %w", err)
	}

	return nil
}

func (r *PacienteRepository) List(ctx context.Context) ([]models.Paciente, error) {
	query := `
		SELECT
			id_paciente,
			nome,
			cpf,
			data_nascimento,
			sexo,
			telefone,
			endereco,
			cartao_sus,
			data_cadastro,
			created_at,
			updated_at
		FROM paciente
		ORDER BY data_cadastro DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar pacientes: %w", err)
	}
	defer rows.Close()

	var pacientes []models.Paciente

	for rows.Next() {
		var paciente models.Paciente

		err := rows.Scan(
			&paciente.ID,
			&paciente.Nome,
			&paciente.CPF,
			&paciente.DataNascimento,
			&paciente.Sexo,
			&paciente.Telefone,
			&paciente.Endereco,
			&paciente.CartaoSUS,
			&paciente.DataCadastro,
			&paciente.CreatedAt,
			&paciente.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao ler paciente: %w", err)
		}

		pacientes = append(pacientes, paciente)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao percorrer pacientes: %w", err)
	}

	return pacientes, nil
}


func (r *PacienteRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Paciente, error) {
	query := `
		SELECT
			id_paciente,
			nome,
			cpf,
			data_nascimento,
			sexo,
			telefone,
			endereco,
			cartao_sus,
			data_cadastro,
			created_at,
			updated_at
		FROM paciente
		WHERE id_paciente = $1
	`

	var paciente models.Paciente

	err := r.db.QueryRow(ctx, query, id).Scan(
		&paciente.ID,
		&paciente.Nome,
		&paciente.CPF,
		&paciente.DataNascimento,
		&paciente.Sexo,
		&paciente.Telefone,
		&paciente.Endereco,
		&paciente.CartaoSUS,
		&paciente.DataCadastro,
		&paciente.CreatedAt,
		&paciente.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar paciente: %w", err)
	}

	return &paciente, nil
}

func (r *PacienteRepository) Update(ctx context.Context, paciente *models.Paciente) error {
	query := `
		UPDATE paciente
		SET
			nome = $1,
			cpf = $2,
			data_nascimento = $3,
			sexo = $4,
			telefone = $5,
			endereco = $6,
			cartao_sus = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE id_paciente = $8
	`

	result, err := r.db.Exec(
		ctx,
		query,
		paciente.Nome,
		paciente.CPF,
		paciente.DataNascimento,
		paciente.Sexo,
		paciente.Telefone,
		paciente.Endereco,
		paciente.CartaoSUS,
		paciente.ID,
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar paciente: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("paciente não encontrado")
	}

	return nil
}

func (r *PacienteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM paciente
		WHERE id_paciente = $1
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar paciente: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("paciente não encontrado")
	}

	return nil
}