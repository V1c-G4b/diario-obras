package storage

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/v1c-g4b/diario-obras/internal/domain/port"
)

var _ port.ObraStorage = (*Storage)(nil)

type Storage struct {
	client     *minio.Client
	bucketName string
}

func NewStorage(client *minio.Client, bucketName string) *Storage {
	return &Storage{
		client:     client,
		bucketName: bucketName,
	}
}

func (s *Storage) Delete(ctx context.Context, fileName string) error {
	return s.client.RemoveObject(ctx, s.bucketName, fileName, minio.RemoveObjectOptions{})
}

func (s *Storage) GetURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	u, err := s.client.PresignedGetObject(ctx, s.bucketName, objectKey, expiry, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (s *Storage) Upload(ctx context.Context, fileName string, file io.Reader, size int64) (string, error) {
	_, err := s.client.PutObject(ctx, s.bucketName, fileName, file, size, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return fileName, nil
}
