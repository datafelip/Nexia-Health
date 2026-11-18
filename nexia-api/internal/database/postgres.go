
package database

import (
	"context"
	"fmt"
	"time"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf(
			"erro ao configurar conexão com PostgreSQL: %w",
			err,
		)
	}

	// Interpreta timestamps sem timezone como
	// horários locais de São Paulo.
	localizacao, err := time.LoadLocation(
		"America/Sao_Paulo",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"erro ao carregar fuso horário: %w",
			err,
		)
	}

	// Configura o codec para cada nova conexão do pool.
	config.AfterConnect = func(
		ctx context.Context,
		conn *pgx.Conn,
	) error {
		conn.TypeMap().RegisterType(&pgtype.Type{
			Name: "timestamp",
			OID:  pgtype.TimestampOID,
			Codec: &pgtype.TimestampCodec{
				ScanLocation: localizacao,
			},
		})

		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf(
			"erro ao criar pool de conexões: %w",
			err,
		)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()

		return nil, fmt.Errorf(
			"erro ao conectar ao PostgreSQL: %w",
			err,
		)
	}

	return pool, nil
}
