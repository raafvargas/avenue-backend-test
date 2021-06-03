package main

import (
	"github.com/avenuesec/file-upload-test/internal/uploader/delivery/http"
	"github.com/avenuesec/file-upload-test/internal/uploader/repository"
	"github.com/avenuesec/file-upload-test/internal/uploader/service"
	"github.com/avenuesec/file-upload-test/pkg/config"
	"github.com/avenuesec/file-upload-test/pkg/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	config, err := config.Load(".env")

	if err != nil {
		panic(err)
	}

	storage, err := storage.NewStorage(config)

	if err != nil {
		panic(err)
	}

	repository := repository.NewInMemoryUploadRepository()
	service := service.NewUploadService(storage, repository)
	controller := http.NewUploadController(service)

	engine := gin.Default()

	controller.RegisterRoutes(engine)

	if err := engine.Run(config.Host); err != nil {
		panic(err)
	}
}
