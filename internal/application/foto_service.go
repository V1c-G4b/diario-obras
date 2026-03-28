package application

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/v1c-g4b/diario-obras/internal/domain/entity"
	"github.com/v1c-g4b/diario-obras/internal/domain/port"
)

type FotoService struct {
	foto    port.FotoRepository
	storage port.ObraStorage
}

func NewFotoService(foto port.FotoRepository, storage port.ObraStorage) *FotoService {
	return &FotoService{foto: foto, storage: storage}
}

func (s *FotoService) Create(ctx context.Context, foto *entity.Foto, fileName string, file io.Reader, size int64) error {
	url, err := s.storage.Upload(ctx, fileName, file, size)
	if err != nil {
		return err
	}
	foto.URLS3 = url

	if err := s.foto.Create(ctx, foto); err != nil {
		if delErr := s.storage.Delete(ctx, fileName); delErr != nil {
			log.Printf("falha ao remover arquivo órfão do storage (%s): %v", fileName, delErr)
		}
		return fmt.Errorf("falha ao salvar foto no banco: %w", err)
	}

	return nil
}

func (s *FotoService) ListByEntrada(ctx context.Context, entradaID uuid.UUID) ([]entity.Foto, error) {
	fotos, err := s.foto.ListByEntrada(ctx, entradaID)

	if err != nil {
		return nil, err
	}

	for i := range fotos {
		url, err := s.storage.GetURL(ctx, fotos[i].URLS3, 15*time.Minute)
		if err != nil {
			return nil, fmt.Errorf("falha ao gerar URL pré-assinada: %w", err)
		}
		fotos[i].URLS3 = url
	}
	return fotos, nil
}

func (s *FotoService) Delete(ctx context.Context, id uuid.UUID) error {
	foto, err := s.foto.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("foto não encontrada: %w", err)
	}

	if err := s.foto.Delete(ctx, id); err != nil {
		return fmt.Errorf("falha ao deletar foto do banco: %w", err)
	}

	if err := s.storage.Delete(ctx, foto.URLS3); err != nil {
		log.Printf("falha ao remover arquivo do storage (%s): %v", foto.URLS3, err)
	}

	return nil
}

func (s *FotoService) DeleteAllByEntrada(ctx context.Context, entradaID uuid.UUID) error {
	fotos, err := s.foto.ListByEntrada(ctx, entradaID)
	if err != nil {
		return err
	}

	for _, foto := range fotos {
		if err := s.foto.Delete(ctx, foto.ID); err != nil {
			log.Printf("falha ao deletar foto %s do banco: %v", foto.ID, err)
			continue
		}
		if err := s.storage.Delete(ctx, foto.URLS3); err != nil {
			log.Printf("falha ao remover arquivo do storage (%s): %v", foto.URLS3, err)
		}
	}

	return nil
}
