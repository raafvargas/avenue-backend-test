package repository

import (
	"context"
	"testing"

	"github.com/avenuesec/file-upload-test/internal/uploader/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UploadRepositoryTestSuite struct {
	suite.Suite

	assert     *assert.Assertions
	repository UploadRepository
}

func TestUploadRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UploadRepositoryTestSuite))
}

func (s *UploadRepositoryTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.repository = NewInMemoryUploadRepository()
}

func (s *UploadRepositoryTestSuite) TestCreateAndGet() {
	ctx := context.Background()

	upload := &domain.Upload{
		FileName: "file-name.txt",
		Location: "bucket/file-name.txt",
	}

	id, err := s.repository.Create(ctx, upload)

	s.assert.NoError(err)

	result, err := s.repository.Get(ctx, id)

	s.assert.NoError(err)
	s.assert.Equal(upload.Location, result.Location)
	s.assert.Equal(upload.FileName, result.FileName)
}

func (s *UploadRepositoryTestSuite) TestGetNotFound() {
	ctx := context.Background()

	_, err := s.repository.Get(ctx, uuid.New().String())

	s.assert.Error(err)
	s.assert.EqualError(err, ErrFileNotFound.Error())
}
