package http

import (
	"io"
	"net/http"

	"github.com/avenuesec/file-upload-test/internal/uploader/service"
	"github.com/avenuesec/file-upload-test/pkg/model"
	"github.com/gin-gonic/gin"
)

type UploadController struct {
	service service.UploadService
}

func NewUploadController(service service.UploadService) *UploadController {
	controller := &UploadController{
		service: service,
	}

	return controller
}

func (ctrl *UploadController) RegisterRoutes(engine *gin.Engine) {
	uploads := engine.Group("uploads")

	uploads.GET(":id", ctrl.Get)
	uploads.POST("", ctrl.Upload)
}

func (ctrl *UploadController) Upload(c *gin.Context) {
	file, err := c.FormFile("file")

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	reader, err := file.Open()

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	upload := &model.Upload{
		File:     reader,
		FileName: file.Filename,
	}

	uploadID, err := ctrl.service.Upload(c.Request.Context(), upload)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Header("Location", uploadID)
	c.Status(http.StatusCreated)
}

func (ctrl *UploadController) Get(c *gin.Context) {
	chunk := &model.ReadChunk{
		UploadID: c.Param("id"),
	}

	if err := c.BindQuery(chunk); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	reader, err := ctrl.service.ReadChunk(c.Request.Context(), chunk)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(c.Writer, reader); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
