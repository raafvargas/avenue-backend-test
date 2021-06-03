package repository

import (
	"context"

	"github.com/avenuesec/file-upload-test/internal/uploader/domain"
	"github.com/stretchr/testify/mock"
)

type MockUploadRepository struct {
	mock.Mock
}

func (m *MockUploadRepository) Get(ctx context.Context, id string) (*domain.Upload, error) {
	args := m.Called(ctx, id)

	upload, _ := args.Get(0).(*domain.Upload)

	return upload, args.Error(1)
}

func (m *MockUploadRepository) Create(ctx context.Context, file *domain.Upload) (string, error) {
	args := m.Called(ctx, file)

	return args.String(0), args.Error(1)
}
