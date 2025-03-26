package http

import (
	"github.com/axlle-com/blog/pkg/app/models"
	"github.com/axlle-com/blog/pkg/file"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	UploadImage(*gin.Context)
	DeleteImage(*gin.Context)
	UploadImages(*gin.Context)
}

func NewFileController(
	service *file.Service,
) Controller {
	return &controller{
		service: service,
	}
}

type controller struct {
	*models.BaseAjax

	service *file.Service
}
