package services

import (
"context"
"fmt"
"time"

"github.com/felipelima/nexia-health/api/internal/models"
"github.com/felipelima/nexia-health/api/internal/repositories"
"github.com/google/uuid"
)

type FilaService struct {
repository *repositories.FilaRepository
}

func NewFilaService(repository *repositories.FilaRepository) *FilaService {
return &FilaService{repository: repository}
}

func statusFilaValido(status models.Status) bool {
switch status {
case models.Aguardando, models.Chamado, models.EmAtendimento, models.Finalizado, models.Canceled:
return true
default:
return false
    }
}

func (s *FilaService) Create(ctx context.Context, fila *models.Fila) error {
if fila.PacienteID == uuid.Nil {
return fmt.Errorf("id_paciente é obrigatório")
    }

if fila.TriagemID == uuid.Nil {
return fmt.Errorf("id_triagem é obrigatório")
    }

if fila.Posicao < 1 {
return fmt.Errorf("posicao deve ser maior ou igual a 1")
    }

if fila.Status == "" {
fila.Status = models.Aguardando
    } else if !statusFilaValido(fila.Status) {
return fmt.Errorf("status inválido")
    }

if fila.DataEntrada.IsZero() {
fila.DataEntrada = time.Now()
    }

return s.repository.Create(ctx, fila)
}

func (s *FilaService) GetByID(ctx context.Context, id uuid.UUID) (*models.Fila, error) {
return s.repository.GetByID(ctx, id)
}

func (s *FilaService) List(ctx context.Context) ([]models.Fila, error) {
return s.repository.List(ctx)
}

// Update atualiza apenas posicao e status (únicos campos que o repository
// realmente altera) e retorna o registro completo após a atualização,
// evitando que id_paciente, id_triagem e data_entrada voltem zerados
// na resposta.
func (s *FilaService) Update(ctx context.Context, fila *models.Fila) (*models.Fila, error) {
if fila.Posicao < 1 {
return nil, fmt.Errorf("posicao deve ser maior ou igual a 1")
    }

if !statusFilaValido(fila.Status) {
return nil, fmt.Errorf("status inválido")
    }

if err := s.repository.Update(ctx, fila); err != nil {
return nil, err
    }

return s.repository.GetByID(ctx, fila.ID)
}

func (s *FilaService) Delete(ctx context.Context, id uuid.UUID) error {
return s.repository.Delete(ctx, id)
}