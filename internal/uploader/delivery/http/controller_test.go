package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/avenuesec/file-upload-test/internal/uploader/service"
	"github.com/avenuesec/file-upload-test/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func init() {
	gin.DefaultWriter = ioutil.Discard
}

type UploadControllerTestSuite struct {
	suite.Suite

	ctx        context.Context
	engine     *gin.Engine
	assert     *assert.Assertions
	service    *service.MockUploadService
	controller *UploadController
}

func TestUploadControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UploadControllerTestSuite))
}

func (s *UploadControllerTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.assert = assert.New(s.T())
	s.engine = gin.New()
	s.service = &service.MockUploadService{}

	s.controller = NewUploadController(s.service)
	s.controller.RegisterRoutes(s.engine)
}

func (s *UploadControllerTestSuite) TearDownTest() {
	s.service.AssertExpectations(s.T())
}

func (s *UploadControllerTestSuite) TestUpload() {
	fileContent := []byte("upload\n")
	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)
	fileWriter, _ := writer.CreateFormFile("file", "file.txt")
	_, err := fileWriter.Write(fileContent)
	s.assert.NoError(err)

	writer.Close()

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/uploads", body)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	uploadID := uuid.New().String()

	s.service.
		On("Upload", s.ctx, mock.AnythingOfType("*model.Upload")).
		Return(uploadID, nil)

	s.engine.ServeHTTP(response, request)

	s.assert.Equal(http.StatusCreated, response.Code)
	s.assert.Equal(uploadID, response.Header().Get("Location"))
}

func (s *UploadControllerTestSuite) TestUploadErr() {
	fileContent := []byte("upload\n")
	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)
	fileWriter, _ := writer.CreateFormFile("file", "file.txt")
	_, err := fileWriter.Write(fileContent)
	s.assert.NoError(err)

	writer.Close()

	response := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/uploads", body)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	s.service.
		On("Upload", s.ctx, mock.AnythingOfType("*model.Upload")).
		Return("", errors.New("some error"))

	s.engine.ServeHTTP(response, request)

	s.assert.Equal(http.StatusInternalServerError, response.Code)
}

func (s *UploadControllerTestSuite) TestUploadWithoutFile() {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/uploads", nil)

	s.engine.ServeHTTP(response, request)

	s.assert.Equal(http.StatusBadRequest, response.Code)
}

func (s *UploadControllerTestSuite) TestGet() {
	offset := int64(0)
	limit := int64(100)
	uploadID := uuid.New().String()

	uri := fmt.Sprintf("/uploads/%s?offset=%d&limit=%d", uploadID, offset, limit)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", uri, nil)

	chunk := &model.ReadChunk{
		Limit:    limit,
		Offset:   offset,
		UploadID: uploadID,
	}

	content := "read chunk\n"

	s.service.On("ReadChunk", s.ctx, chunk).
		Return(bytes.NewBufferString(content), nil)

	s.engine.ServeHTTP(response, request)

	s.assert.Equal(http.StatusOK, response.Code)

	resBody, _ := ioutil.ReadAll(response.Body)

	s.assert.Equal(content, string(resBody))
}

func (s *UploadControllerTestSuite) TestGetError() {
	offset := int64(0)
	limit := int64(100)
	uploadID := uuid.New().String()

	uri := fmt.Sprintf("/uploads/%s?offset=%d&limit=%d", uploadID, offset, limit)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", uri, nil)

	chunk := &model.ReadChunk{
		Limit:    limit,
		Offset:   offset,
		UploadID: uploadID,
	}

	s.service.On("ReadChunk", s.ctx, chunk).
		Return(nil, errors.New("some error"))

	s.engine.ServeHTTP(response, request)

	s.assert.Equal(http.StatusInternalServerError, response.Code)
}

func (s *UploadControllerTestSuite) TestGetWithoutChunk() {
	uploadID := uuid.New().String()

	uri := fmt.Sprintf("/uploads/%s", uploadID)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", uri, nil)

	s.engine.ServeHTTP(response, request)

	s.assert.Equal(http.StatusBadRequest, response.Code)
}
