package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/v1c-g4b/diario-obras/internal/domain/entity"
	"github.com/v1c-g4b/diario-obras/internal/domain/port"
)

type EntradaService struct {
	entrada     port.EntradaRepository
	obra        port.ObraRepository
	responsavel port.ResponsavelRepository
	foto        *FotoService
}

func NewEntradaService(entrada port.EntradaRepository, obra port.ObraRepository, foto *FotoService, responsavel port.ResponsavelRepository) *EntradaService {
	return &EntradaService{entrada: entrada, obra: obra, foto: foto, responsavel: responsavel}
}

func (s *EntradaService) Create(ctx context.Context, entrada *entity.Entrada, obraId uuid.UUID) error {

	obra, err := s.obra.FindByID(ctx, obraId)

	if err != nil {
		return err
	}

	if obra == nil {
		return errors.New("Obra não encontrada")
	}

	responsavel, err := s.responsavel.FindByID(ctx, entrada.ResponsavelID)

	if err != nil {
		return err
	}

	if responsavel == nil {
		return errors.New("Responsavel não encontrado")
	}

	entrada.ObraID = obraId
	entrada.Responsavel = *responsavel
	return s.entrada.Create(ctx, entrada)
}

func (s *EntradaService) FindByID(ctx context.Context, id uuid.UUID) (*entity.Entrada, error) {
	return s.entrada.FindByID(ctx, id)
}

func (s *EntradaService) ListByObra(ctx context.Context, obraId uuid.UUID) ([]entity.Entrada, error) {
	return s.entrada.ListByObra(ctx, obraId)
}

func (s *EntradaService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.foto.DeleteAllByEntrada(ctx, id); err != nil {
		return fmt.Errorf("falha ao remover fotos da entrada: %w", err)
	}
	return s.entrada.Delete(ctx, id)
}
