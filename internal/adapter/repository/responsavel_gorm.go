package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/v1c-g4b/diario-obras/internal/domain/entity"
	"github.com/v1c-g4b/diario-obras/internal/domain/port"
	"gorm.io/gorm"
)

type ResponsavelGormRepository struct {
	db *gorm.DB
}

func NewResponsavelGormRepository(db *gorm.DB) *ResponsavelGormRepository {
	return &ResponsavelGormRepository{db: db}
}

var _ port.ResponsavelRepository = (*ResponsavelGormRepository)(nil)

func (r *ResponsavelGormRepository) Create(ctx context.Context, responsavel *entity.Responsavel) error {
	return r.db.WithContext(ctx).Create(responsavel).Error
}

func (r *ResponsavelGormRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Responsavel, error) {
	var responsavel entity.Responsavel
	err := r.db.WithContext(ctx).First(&responsavel, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &responsavel, nil
}

func (r *ResponsavelGormRepository) List(ctx context.Context) ([]entity.Responsavel, error) {
	var responsaveis []entity.Responsavel
	err := r.db.WithContext(ctx).Find(&responsaveis).Error
	return responsaveis, err
}

func (r *ResponsavelGormRepository) Update(ctx context.Context, responsavel *entity.Responsavel) error {
	return r.db.WithContext(ctx).Save(responsavel).Error
}

func (r *ResponsavelGormRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Responsavel{}, "id = ?", id).Error
}
