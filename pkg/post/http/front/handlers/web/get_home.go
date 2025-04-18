package web

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/axlle-com/blog/pkg/post/models"
)

func (c *postController) GetHome(ctx *gin.Context) {
	post, err := c.postService.GetByParam("is_main", true)
	if err != nil || post == nil {
		logger.Error(err)
		post = &models.Post{}
	}

	ctx.HTML(
		http.StatusOK,
		"index",
		gin.H{
			"title": "Home Page",
			"post":  post,
		},
	)
}
