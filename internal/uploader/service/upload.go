package service

import (
	"context"
	"fmt"
	"io"

	"github.com/avenuesec/file-upload-test/internal/uploader/domain"
	"github.com/avenuesec/file-upload-test/internal/uploader/repository"
	"github.com/avenuesec/file-upload-test/pkg/model"
	"github.com/avenuesec/file-upload-test/pkg/storage"
	"github.com/google/uuid"
)

type UploadService interface {
	Upload(ctx context.Context, upload *model.Upload) (string, error)
	ReadChunk(ctx context.Context, chunk *model.ReadChunk) (io.Reader, error)
}

type uploadService struct {
	storage    storage.Storage
	repository repository.UploadRepository
}

func NewUploadService(storage storage.Storage, repo repository.UploadRepository) UploadService {
	return &uploadService{
		repository: repo,
		storage:    storage,
	}
}

func (s *uploadService) Upload(ctx context.Context, upload *model.Upload) (string, error) {
	// TODO: I should inject an uuid provider here. cannot test this way
	location := fmt.Sprintf("%s-%s", uuid.New().String(), upload.FileName)

	if err := s.storage.Copy(ctx, location, upload.File); err != nil {
		return "", err
	}

	id, err := s.repository.Create(ctx, &domain.Upload{
		Location: location,
		FileName: upload.FileName,
	})

	return id, err
}

func (s *uploadService) ReadChunk(ctx context.Context, chunk *model.ReadChunk) (io.Reader, error) {
	upload, err := s.repository.Get(ctx, chunk.UploadID)

	if err != nil {
		return nil, err
	}

	reader, err := s.storage.Read(ctx, upload.Location, chunk.Offset, chunk.Limit)

	return reader, err
}
