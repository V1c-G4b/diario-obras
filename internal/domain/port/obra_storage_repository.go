package port

import (
	"context"
	"io"
	"time"
)

type ObraStorage interface {
	Upload(ctx context.Context, fileName string, file io.Reader, size int64) (string, error)
	GetURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error)
	Delete(ctx context.Context, fileName string) error
}
