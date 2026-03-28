package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/v1c-g4b/diario-obras/internal/domain/entity"
)

type FotoRepository interface {
	Create(ctx context.Context, foto *entity.Foto) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Foto, error)
	ListByEntrada(ctx context.Context, entradaID uuid.UUID) ([]entity.Foto, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
