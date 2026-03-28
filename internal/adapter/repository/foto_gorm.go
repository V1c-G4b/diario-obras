package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/v1c-g4b/diario-obras/internal/domain/entity"
	"github.com/v1c-g4b/diario-obras/internal/domain/port"
	"gorm.io/gorm"
)

var _ port.FotoRepository = (*FotoGormRepository)(nil)

type FotoGormRepository struct {
	db *gorm.DB
}

func (r *FotoGormRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Foto, error) {
	var foto entity.Foto
	err := r.db.WithContext(ctx).First(&foto, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &foto, nil
}

func (r *FotoGormRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Foto{}, "id = ?", id).Error
}

func (r *FotoGormRepository) ListByEntrada(ctx context.Context, entradaID uuid.UUID) ([]entity.Foto, error) {
	var fotos []entity.Foto
	err := r.db.WithContext(ctx).Where("entrada_id = ?", entradaID).Find(&fotos).Error
	return fotos, err
}

func NewFotoGormRepository(db *gorm.DB) *FotoGormRepository {
	return &FotoGormRepository{db: db}
}

func (r *FotoGormRepository) Create(ctx context.Context, foto *entity.Foto) error {
	return r.db.WithContext(ctx).Create(foto).Error
}
