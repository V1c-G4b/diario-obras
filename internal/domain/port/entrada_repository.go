package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/v1c-g4b/diario-obras/internal/domain/entity"
)

type EntradaRepository interface {
	Create(ctx context.Context, entrada *entity.Entrada) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Entrada, error)
	ListByObra(ctx context.Context, obraId uuid.UUID) ([]entity.Entrada, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
