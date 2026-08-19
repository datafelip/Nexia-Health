package main

import (
	"log"
	"net/http"

	"github.com/felipelima/nexia-health/api/configs"
	"github.com/felipelima/nexia-health/api/internal/database"
	"github.com/felipelima/nexia-health/api/internal/handlers"
	"github.com/felipelima/nexia-health/api/internal/repositories"
	"github.com/felipelima/nexia-health/api/internal/routes"
	"github.com/felipelima/nexia-health/api/internal/services"
)

func main() {
	config, err := configs.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Repositories
	pacienteRepository := repositories.NewPacienteRepository(db)

	// Services
	pacienteService := services.NewPacienteService(pacienteRepository)

	// Handlers
	pacienteHandler := handlers.NewPacienteHandler(pacienteService)

	// Routes
	router := routes.Setup(pacienteHandler)

	log.Println("Nexia Health API")
	log.Println("PostgreSQL conectado com sucesso!")
	log.Println("Servidor iniciado em http://localhost:8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}