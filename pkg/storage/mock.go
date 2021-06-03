package storage

import (
	"context"
	"io"

	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Copy(ctx context.Context, fileName string, reader io.Reader) error {
	return m.Called(ctx, fileName, reader).Error(0)
}

func (m *MockStorage) Read(ctx context.Context, fileName string, offset, limit int64) (io.Reader, error) {
	args := m.Called(ctx, fileName, offset, limit)

	reader, _ := args.Get(0).(io.Reader)

	return reader, args.Error(1)
}
