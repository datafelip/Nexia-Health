package main

import (
	"github.com/felipelima/nexia-health/api/configs"
	"github.com/felipelima/nexia-health/api/internal/database"
	"github.com/felipelima/nexia-health/api/internal/handlers"
	"github.com/felipelima/nexia-health/api/internal/repositories"
	"github.com/felipelima/nexia-health/api/internal/routes"
	"github.com/felipelima/nexia-health/api/internal/services"
	"log"
	"net/http"
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
	triagemRepository := repositories.NewTriagemRepository(db)
	filaRepository := repositories.NewFilaRepository(db)
	atendimentoRepository := repositories.NewAtendimentoRepository(db)
	agendamentoRepository := repositories.NewAgendamentoRepository(db)
	especialidadeRepository := repositories.NewEspecialidadeRepository(db)
	setorRepository := repositories.NewSetorRepository(db)
	usuarioRepository := repositories.NewUsuarioRepository(db)

	// Services
	pacienteService := services.NewPacienteService(pacienteRepository)
	triagemService := services.NewTriagemService(triagemRepository)
	filaService := services.NewFilaService(filaRepository)
	atendimentoService := services.NewAtendimentoService(atendimentoRepository)
	agendamentoService := services.NewAgendamentoService(agendamentoRepository)
	especialidadeService := services.NewEspecialidadeService(especialidadeRepository)
	setorService := services.NewSetorService(setorRepository)
	usuarioService := services.NewUsuarioService(usuarioRepository)

	// Handlers
	pacienteHandler := handlers.NewPacienteHandler(pacienteService)
	triagemHandler := handlers.NewTriagemHandler(triagemService)
	filaHandler := handlers.NewFilaHandler(filaService)
	atendimentoHandler := handlers.NewAtendimentoHandler(atendimentoService)
	agendamentoHandler := handlers.NewAgendamentoHandler(agendamentoService)
	especialidadeHandler := handlers.NewEspecialidadeHandler(especialidadeService)
	setorHandler := handlers.NewSetorHandler(setorService)
	usuarioHandler := handlers.NewUsuarioHandler(usuarioService)

	// Routes
	router := routes.Setup(
		pacienteHandler,
		triagemHandler,
		filaHandler,
		atendimentoHandler,
		agendamentoHandler,
		especialidadeHandler,
		setorHandler,
		usuarioHandler,
	)

	log.Println("Nexia Health API")
	log.Println("PostgreSQL conectado com sucesso!")
	log.Println("Servidor iniciado em http://localhost:8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
