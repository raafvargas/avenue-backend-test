package service

import (
	"context"
	"io"

	"github.com/avenuesec/file-upload-test/pkg/model"
	"github.com/stretchr/testify/mock"
)

type MockUploadService struct {
	mock.Mock
}

func (m *MockUploadService) Upload(ctx context.Context, upload *model.Upload) (string, error) {
	args := m.Called(ctx, upload)

	return args.String(0), args.Error(1)
}

func (m *MockUploadService) ReadChunk(ctx context.Context, chunk *model.ReadChunk) (io.Reader, error) {
	args := m.Called(ctx, chunk)

	reader, _ := args.Get(0).(io.Reader)

	return reader, args.Error(1)
}
