package service

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	"github.com/avenuesec/file-upload-test/internal/uploader/domain"
	"github.com/avenuesec/file-upload-test/internal/uploader/repository"
	"github.com/avenuesec/file-upload-test/pkg/model"
	"github.com/avenuesec/file-upload-test/pkg/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UploadServiceTestSuite struct {
	suite.Suite

	ctx        context.Context
	assert     *assert.Assertions
	storage    *storage.MockStorage
	repository *repository.MockUploadRepository
	service    UploadService
}

func TestUploadServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UploadServiceTestSuite))
}

func (s *UploadServiceTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.assert = assert.New(s.T())
	s.storage = &storage.MockStorage{}
	s.repository = &repository.MockUploadRepository{}
	s.service = NewUploadService(
		s.storage,
		s.repository,
	)
}

func (s *UploadServiceTestSuite) TearDownTest() {
	s.storage.AssertExpectations(s.T())
	s.repository.AssertExpectations(s.T())
}

func (s *UploadServiceTestSuite) TestUpload() {
	buffer := bytes.NewBufferString("upload test")

	upload := &model.Upload{
		FileName: "file.txt",
		File:     buffer,
	}

	id := uuid.New().String()

	s.storage.
		On("Copy", s.ctx, mock.AnythingOfType("string"), buffer).
		Return(nil)

	s.repository.
		On("Create", s.ctx, mock.AnythingOfType("*domain.Upload")).
		Return(id, nil)

	result, err := s.service.Upload(s.ctx, upload)

	s.assert.NoError(err)
	s.assert.Equal(id, result)
}

func (s *UploadServiceTestSuite) TestReadChunk() {
	location := "bucket/file.txt"

	chunk := &model.ReadChunk{
		UploadID: uuid.New().String(),
		Offset:   0,
		Limit:    100,
	}

	buffer := bytes.NewBufferString("upload test")

	s.repository.
		On("Get", s.ctx, chunk.UploadID).
		Return(&domain.Upload{
			ID:       chunk.UploadID,
			Location: location,
		}, nil)

	s.storage.
		On("Read", s.ctx, location, chunk.Offset, chunk.Limit).
		Return(buffer, nil)

	reader, err := s.service.ReadChunk(s.ctx, chunk)

	s.assert.NoError(err)

	content, _ := ioutil.ReadAll(reader)

	s.assert.Equal("upload test", string(content))
}

func (s *UploadServiceTestSuite) TestReadChunkUploadNotFound() {

	chunk := &model.ReadChunk{
		UploadID: uuid.New().String(),
		Offset:   0,
		Limit:    100,
	}

	s.repository.
		On("Get", s.ctx, chunk.UploadID).
		Return(nil, repository.ErrFileNotFound)

	_, err := s.service.ReadChunk(s.ctx, chunk)

	s.assert.Error(err)
	s.assert.EqualError(err, repository.ErrFileNotFound.Error())
}
