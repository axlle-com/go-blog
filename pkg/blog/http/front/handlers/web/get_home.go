package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/blog/models"
)

func (c *postController) GetHome(ctx *gin.Context) {
	post, err := c.postService.GetByParam("is_main", true)
	if err != nil || post == nil {
		logger.Debugf("[PostController][GetHome] Error: %v", err)
		post = &models.Post{}
	}

	ctx.HTML(
		http.StatusOK,
		c.view.View(nil),
		gin.H{
			"title": "Home Page",
			"post":  post,
		},
	)
}
