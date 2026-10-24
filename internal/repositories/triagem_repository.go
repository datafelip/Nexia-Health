package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/felipelima/nexia-health/api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var ErrTriagemNotFound = errors.New("triagem não encontrada")

var ErrTriagemNaoAguardando = errors.New(
	"somente pacientes aguardando podem ser reclassificados",
)

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

func (r *TriagemRepository) CreateComFila(
	ctx context.Context,
	triagem *models.Triagem,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	// Serializa as operações de ordenação que usam este fluxo.
	_, err = tx.Exec(ctx,
		"SELECT pg_advisory_xact_lock(874321)",
	)
	if err != nil {
		return fmt.Errorf("erro ao bloquear fila: %w", err)
	}

	// Evita inserir um paciente que já está em uma fila ativa.
	var jaExiste bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM fila
			WHERE id_paciente = $1
			  AND status IN (
			    'Aguardando', 'Chamado', 'Em Atendimento'
			  )
		)
	`, triagem.PacienteID).Scan(&jaExiste)
	if err != nil {
		return fmt.Errorf("erro ao verificar fila: %w", err)
	}
	if jaExiste {
		return fmt.Errorf("paciente já possui entrada ativa na fila")
	}

	err = tx.QueryRow(ctx, `
		INSERT INTO triagem (
			id_paciente, id_profissional,
			classificacao_risco, pressao_arterial,
			temperatura, frequencia_cardiaca, observacoes
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id_triagem, data_triagem,
		          created_at, updated_at
	`,
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
		return fmt.Errorf("erro ao salvar triagem: %w", err)
	}

	// A posição provisória será recalculada ainda nesta transação.
	var dataEntrada time.Time
	err = tx.QueryRow(ctx, `
		INSERT INTO fila (
			id_paciente, id_triagem,
			posicao, status
		)
		VALUES ($1, $2, 1, 'Aguardando')
		RETURNING data_entrada
	`, triagem.PacienteID, triagem.ID).Scan(&dataEntrada)
	if err != nil {
		return fmt.Errorf("erro ao inserir na fila: %w", err)
	}

	// Reordena exclusivamente quem ainda está aguardando.
	_, err = tx.Exec(ctx, `
		WITH ordenacao AS (
			SELECT
				f.id_fila,
				ROW_NUMBER() OVER (
					ORDER BY
						CASE t.classificacao_risco
							WHEN 'Vermelho' THEN 1
							WHEN 'Amarelo' THEN 2
							WHEN 'Verde' THEN 3
						END,
						f.data_entrada ASC,
						f.id_fila ASC
				)::integer AS nova_posicao
			FROM fila f
			JOIN triagem t ON t.id_triagem = f.id_triagem
			WHERE f.status = 'Aguardando'
		)
		UPDATE fila f
		SET posicao = o.nova_posicao,
		    updated_at = NOW()
		FROM ordenacao o
		WHERE f.id_fila = o.id_fila
		  AND f.posicao IS DISTINCT FROM o.nova_posicao
	`)
	if err != nil {
		return fmt.Errorf("erro ao reorganizar fila: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("erro ao confirmar transação: %w", err)
	}

	return nil
}

func (r *TriagemRepository) ReclassificarComFila(
	ctx context.Context,
	id uuid.UUID,
	novaClassificacao models.ClassificacaoRisco,
) (*models.Triagem, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	// Utiliza o mesmo bloqueio das demais operações da fila.
	_, err = tx.Exec(ctx,
		"SELECT pg_advisory_xact_lock(874321)",
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao bloquear fila: %w", err)
	}

	// Verifica se a triagem existe.
	var existe bool

	err = tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM triagem
			WHERE id_triagem = $1
		)
	`, id).Scan(&existe)

	if err != nil {
		return nil, fmt.Errorf("erro ao verificar triagem: %w", err)
	}

	if !existe {
		return nil, ErrTriagemNotFound
	}

	// A reclassificação só é permitida para quem aguarda.
	var filaID uuid.UUID

	err = tx.QueryRow(ctx, `
		SELECT id_fila
		FROM fila
		WHERE id_triagem = $1
		  AND status = 'Aguardando'
		FOR UPDATE
	`, id).Scan(&filaID)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTriagemNaoAguardando
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao consultar fila: %w", err)
	}

	// Atualiza a classificação de risco.
	_, err = tx.Exec(ctx, `
		UPDATE triagem
		SET classificacao_risco = $1,
		    updated_at = NOW()
		WHERE id_triagem = $2
	`, novaClassificacao, id)

	if err != nil {
		return nil, fmt.Errorf("erro ao reclassificar triagem: %w", err)
	}

	// Recalcula a posição de todos os pacientes aguardando.
	_, err = tx.Exec(ctx, `
		WITH ordenacao AS (
			SELECT
				f.id_fila,
				ROW_NUMBER() OVER (
					ORDER BY
						CASE t.classificacao_risco
							WHEN 'Vermelho' THEN 1
							WHEN 'Amarelo' THEN 2
							WHEN 'Verde' THEN 3
						END,
						f.data_entrada ASC,
						f.id_fila ASC
				)::integer AS nova_posicao
			FROM fila f
			JOIN triagem t
			  ON t.id_triagem = f.id_triagem
			WHERE f.status = 'Aguardando'
		)
		UPDATE fila f
		SET posicao = o.nova_posicao,
		    updated_at = NOW()
		FROM ordenacao o
		WHERE f.id_fila = o.id_fila
		  AND f.posicao IS DISTINCT FROM o.nova_posicao
	`)

	if err != nil {
		return nil, fmt.Errorf("erro ao reorganizar fila: %w", err)
	}

	// Busca a triagem atualizada dentro da transação.
	var triagem models.Triagem

	err = tx.QueryRow(ctx, `
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
	`, id).Scan(
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
		return nil, fmt.Errorf("erro ao buscar triagem atualizada: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("erro ao confirmar reclassificação: %w", err)
	}

	return &triagem, nil
}
