package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/avenuesec/file-upload-test/internal/uploader/domain"
	"github.com/google/uuid"
)

var (
	ErrFileNotFound = errors.New("file not found")
)

type UploadRepository interface {
	Get(ctx context.Context, id string) (*domain.Upload, error)
	Create(ctx context.Context, file *domain.Upload) (string, error)
}

type uploadRepository struct {
	store sync.Map
}

func NewInMemoryUploadRepository() UploadRepository {
	return &uploadRepository{
		store: sync.Map{},
	}
}

func (r *uploadRepository) Get(ctx context.Context, id string) (*domain.Upload, error) {
	upload, exists := r.store.Load(id)

	if !exists {
		return nil, ErrFileNotFound
	}

	return upload.(*domain.Upload), nil
}

func (r *uploadRepository) Create(ctx context.Context, file *domain.Upload) (string, error) {
	file.ID = uuid.New().String()

	r.store.Store(file.ID, file)

	return file.ID, nil
}
