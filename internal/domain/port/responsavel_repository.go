package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/v1c-g4b/diario-obras/internal/domain/entity"
)

type ResponsavelRepository interface {
	Create(ctx context.Context, responsavel *entity.Responsavel) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Responsavel, error)
	List(ctx context.Context) ([]entity.Responsavel, error)
	Update(ctx context.Context, responsavel *entity.Responsavel) error
	Delete(ctx context.Context, id uuid.UUID) error
}
