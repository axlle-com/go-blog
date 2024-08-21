package web

import (
	"github.com/axlle-com/blog/pkg/common/models"
	. "github.com/axlle-com/blog/pkg/post/service"
	"github.com/gin-gonic/gin"
)

type WebController interface {
	GetPost(*gin.Context)
	getPosts(*gin.Context, Container)
	GetPostsHandler() gin.HandlerFunc
	CreatePost(*gin.Context)
}

func NewWebController(r *gin.Engine) WebController {
	return &webController{engine: r}
}

type webController struct {
	*models.BaseAjax
	engine *gin.Engine
}

func (c *webController) GetPostsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		container := NewContainer()
		c.getPosts(ctx, container)
	}
}
