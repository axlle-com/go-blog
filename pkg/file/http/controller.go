package http

import (
	"github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/file/service"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	UploadImage(*gin.Context)
	DeleteImage(*gin.Context)
	UploadImages(*gin.Context)
}

func NewFileController(
	uploadService *service.UploadService,
	fileService *service.FileService,
) Controller {
	return &controller{
		uploadService: uploadService,
		fileService:   fileService,
	}
}

type controller struct {
	*models.BaseAjax

	uploadService *service.UploadService
	fileService   *service.FileService
}
