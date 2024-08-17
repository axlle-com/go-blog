package web

import (
	"github.com/axlle-com/blog/pkg/common/models"
	"github.com/gin-gonic/gin"
)

type WebController interface {
	GetPost(*gin.Context)
	GetPosts(*gin.Context)
	CreatePost(*gin.Context)
}

func NewWebController(r *gin.Engine) WebController {
	return &webController{engine: r}
}

type webController struct {
	*models.BaseAjax
	engine *gin.Engine
}
