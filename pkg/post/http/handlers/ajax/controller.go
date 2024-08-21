package ajax

import (
	"github.com/axlle-com/blog/pkg/common/models"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	UpdatePostHandler() gin.HandlerFunc
	CreatePost(*gin.Context)
	DeletePost(*gin.Context)
	FilterPostHandler() gin.HandlerFunc
	DeletePostImageHandler() gin.HandlerFunc
}

func NewController(r *gin.Engine) Controller {
	return &controller{engine: r}
}

type controller struct {
	*models.BaseAjax
	engine *gin.Engine
}
