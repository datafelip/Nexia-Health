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

var ErrPacienteNotFound = errors.New("Paciente não encontrado")

var ErrPacienteFilaAtiva = errors.New(
	"não é possível excluir um paciente que está na fila de atendimento",
)

var ErrPacienteAtendimentoAtivo = errors.New(
	"não é possível excluir um paciente com atendimento em andamento",
)

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
		WHERE excluido_em IS NULL
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
			AND excluido_em IS NULL
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
			AND excluido_em IS NULL
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
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	// Bloqueia o registro do paciente durante a operação.
	var pacienteID uuid.UUID

	err = tx.QueryRow(ctx, `
		SELECT id_paciente
		FROM paciente
		WHERE id_paciente = $1
		  AND excluido_em IS NULL
		FOR UPDATE
	`, id).Scan(&pacienteID)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrPacienteNotFound
	}
	if err != nil {
		return fmt.Errorf("erro ao consultar paciente: %w", err)
	}

	// Verifica se o paciente está em uma fila ativa.
	var possuiFilaAtiva bool

	err = tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM fila
			WHERE id_paciente = $1
			  AND status IN (
				  'Aguardando',
				  'Chamado',
				  'Em Atendimento'
			  )
		)
	`, id).Scan(&possuiFilaAtiva)

	if err != nil {
		return fmt.Errorf("erro ao verificar fila: %w", err)
	}

	if possuiFilaAtiva {
		return ErrPacienteFilaAtiva
	}

	// Verifica se existe atendimento em andamento.
	var possuiAtendimentoAtivo bool

	err = tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM atendimento
			WHERE id_paciente = $1
			  AND data_inicio IS NOT NULL
			  AND data_fim IS NULL
		)
	`, id).Scan(&possuiAtendimentoAtivo)

	if err != nil {
		return fmt.Errorf("erro ao verificar atendimento: %w", err)
	}

	if possuiAtendimentoAtivo {
		return ErrPacienteAtendimentoAtivo
	}

	// Exclusão lógica.
	result, err := tx.Exec(ctx, `
		UPDATE paciente
		SET
			excluido_em = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
		WHERE id_paciente = $1
		  AND excluido_em IS NULL
	`, id)

	if err != nil {
		return fmt.Errorf("erro ao inativar paciente: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrPacienteNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("erro ao confirmar exclusão: %w", err)
	}

	return nil
}
