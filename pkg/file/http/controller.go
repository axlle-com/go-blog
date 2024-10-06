package http

import (
	"github.com/axlle-com/blog/pkg/common/models"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	UploadImage(*gin.Context)
	DeleteImage(*gin.Context)
	UploadImages(*gin.Context)
}

func NewController(r *gin.Engine) Controller {
	return &controller{engine: r}
}

type controller struct {
	*models.BaseAjax
	engine *gin.Engine
}
