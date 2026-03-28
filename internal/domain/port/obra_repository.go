package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/v1c-g4b/diario-obras/internal/domain/entity"
)

type ObraRepository interface {
	Create(ctx context.Context, obra *entity.Obra) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Obra, error)
	List(ctx context.Context) ([]entity.Obra, error)
	Update(ctx context.Context, obra *entity.Obra) error
	Delete(ctx context.Context, id uuid.UUID) error
}
