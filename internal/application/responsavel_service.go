package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/v1c-g4b/diario-obras/internal/domain/entity"
	"github.com/v1c-g4b/diario-obras/internal/domain/port"
)

type ResponsavelService struct {
	repo port.ResponsavelRepository
}

func NewResponsavelService(repo port.ResponsavelRepository) *ResponsavelService {
	return &ResponsavelService{repo: repo}
}

func (s *ResponsavelService) Create(ctx context.Context, responsavel *entity.Responsavel) error {
	return s.repo.Create(ctx, responsavel)
}

func (s *ResponsavelService) FindByID(ctx context.Context, id uuid.UUID) (*entity.Responsavel, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ResponsavelService) List(ctx context.Context) ([]entity.Responsavel, error) {
	return s.repo.List(ctx)
}

func (s *ResponsavelService) Update(ctx context.Context, responsavel *entity.Responsavel) error {
	return s.repo.Update(ctx, responsavel)
}

func (s *ResponsavelService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
