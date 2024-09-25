package ajax

import (
	"github.com/axlle-com/blog/pkg/common/models"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	UpdatePost(*gin.Context)
	CreatePost(*gin.Context)
	DeletePost(*gin.Context)
	DeletePostImage(*gin.Context)
	FilterPosts(*gin.Context)
}

func NewController(r *gin.Engine) Controller {
	return &controller{engine: r}
}

type controller struct {
	*models.BaseAjax
	engine *gin.Engine
}
