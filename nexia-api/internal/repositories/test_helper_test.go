package repositories

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/felipelima/nexia-health/api/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	// Os testes precisam de uma conexão própria.
	// Nunca utiliza automaticamente o banco da aplicação.
	databaseURL := strings.TrimSpace(
		os.Getenv("TEST_DATABASE_URL"),
	)

	if databaseURL == "" {
		t.Fatal(
			"TEST_DATABASE_URL não configurada. " +
				"Configure um banco exclusivo para testes.",
		)
	}

	// A URL deve apontar explicitamente para nexia_test_db.
	// Isso impede o uso acidental de nexia_db.
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		t.Fatalf(
			"URL do banco de testes inválida: %v",
			err,
		)
	}

	if config.ConnConfig.Database != "nexia_test_db" {
		t.Fatalf(
			"banco de testes inválido: esperado nexia_test_db, recebido %q",
			config.ConnConfig.Database,
		)
	}

	db, err := database.Connect(databaseURL)
	if err != nil {
		t.Fatalf(
			"erro ao conectar ao banco de testes: %v",
			err,
		)
	}

	// Confirma também o banco efetivamente conectado.
	var bancoAtual string

	err = db.QueryRow(
		context.Background(),
		"SELECT current_database()",
	).Scan(&bancoAtual)

	if err != nil {
		db.Close()
		t.Fatalf(
			"erro ao verificar banco de testes: %v",
			err,
		)
	}

	if bancoAtual != "nexia_test_db" {
		db.Close()
		t.Fatalf(
			"conexão recusada: banco atual é %q",
			bancoAtual,
		)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
