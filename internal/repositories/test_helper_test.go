package repositories

import (
	"os"
	"testing"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/felipelima/nexia-health/api/configs"
	"github.com/felipelima/nexia-health/api/internal/database"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	config, err := configs.Load()
	if err != nil {
		t.Fatalf("erro ao carregar configurações: %v", err)
	}

	db, err := database.Connect(config.DatabaseURL)
	if err != nil {
		t.Fatalf("erro ao conectar ao banco de testes: %v", err)
	}

	return db
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}