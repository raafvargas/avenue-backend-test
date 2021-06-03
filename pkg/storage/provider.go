package storage

import (
	"context"
	"errors"
	"io"

	"github.com/avenuesec/file-upload-test/pkg/config"
)

var (
	ErrUnknownProvider = errors.New("unknown storage provider")
)

type Storage interface {
	Copy(ctx context.Context, fileName string, reader io.Reader) error
	Read(ctx context.Context, fileName string, offset, limit int64) (io.Reader, error)
}

func NewStorage(config *config.Configuration) (Storage, error) {
	switch config.StorageBackend {
	case "gcs":
		return NewGCSStorage(config)
	case "s3":
		return NewS3Storage(config)
	}

	return nil, ErrUnknownProvider
}
