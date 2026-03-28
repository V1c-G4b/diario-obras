package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/v1c-g4b/diario-obras/internal/domain/entity"
	"github.com/v1c-g4b/diario-obras/internal/domain/port"
)

type ObraService struct {
	repo port.ObraRepository
}

func NewObraService(repo port.ObraRepository) *ObraService {
	return &ObraService{repo: repo}
}

func (s *ObraService) Create(ctx context.Context, obra *entity.Obra) error {
	return s.repo.Create(ctx, obra)
}

func (s *ObraService) FindByID(ctx context.Context, id uuid.UUID) (*entity.Obra, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ObraService) List(ctx context.Context) ([]entity.Obra, error) {
	return s.repo.List(ctx)
}

func (s *ObraService) Update(ctx context.Context, obra *entity.Obra) error {
	return s.repo.Update(ctx, obra)
}

func (s *ObraService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
