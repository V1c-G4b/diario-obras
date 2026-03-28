package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/v1c-g4b/diario-obras/internal/domain/entity"
	"github.com/v1c-g4b/diario-obras/internal/domain/port"
	"gorm.io/gorm"
)

var _ port.ObraRepository = (*ObraGormRepository)(nil)

type ObraGormRepository struct {
	db *gorm.DB
}

func NewObraGormRepository(db *gorm.DB) *ObraGormRepository {
	return &ObraGormRepository{db: db}
}

func (r *ObraGormRepository) Create(ctx context.Context, obra *entity.Obra) error {
	return r.db.WithContext(ctx).Create(obra).Error
}

func (r *ObraGormRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Obra, error) {
	var obra entity.Obra
	err := r.db.WithContext(ctx).Preload("Entradas").First(&obra, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &obra, nil
}

func (r *ObraGormRepository) List(ctx context.Context) ([]entity.Obra, error) {
	var obras []entity.Obra
	err := r.db.WithContext(ctx).Find(&obras).Error
	return obras, err
}

func (r *ObraGormRepository) Update(ctx context.Context, obra *entity.Obra) error {
	return r.db.WithContext(ctx).Save(obra).Error
}

func (r *ObraGormRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Obra{}, "id = ?", id).Error
}
