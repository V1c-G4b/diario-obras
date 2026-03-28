package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/v1c-g4b/diario-obras/internal/domain/entity"
	"github.com/v1c-g4b/diario-obras/internal/domain/port"
	"gorm.io/gorm"
)

var _ port.EntradaRepository = (*EntradaGormRepository)(nil)

type EntradaGormRepository struct {
	db *gorm.DB
}

func NewEntradaGormRepository(db *gorm.DB) *EntradaGormRepository {
	return &EntradaGormRepository{db: db}
}

func (r *EntradaGormRepository) Create(ctx context.Context, entrada *entity.Entrada) error {
	return r.db.WithContext(ctx).Create(entrada).Error
}

func (r *EntradaGormRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Entrada, error) {
	var entrada entity.Entrada
	err := r.db.WithContext(ctx).Preload("Fotos").Preload("Responsavel").First(&entrada, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entrada, nil
}

func (r *EntradaGormRepository) ListByObra(ctx context.Context, obraId uuid.UUID) ([]entity.Entrada, error) {
	var entradas []entity.Entrada
	err := r.db.WithContext(ctx).Preload("Fotos").Where("obra_id = ?", obraId).Find(&entradas).Error
	return entradas, err
}

func (r *EntradaGormRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Entrada{}, "id = ?", id).Error
}
