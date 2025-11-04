package web

import (
	"encoding/json"
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/dto"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/gin-gonic/gin"
)

func (c *postController) GetHome(ctx *gin.Context) {
	post := &models.Post{}
	var blocks []dto.InfoBlock

	post, err := c.postService.FindByParam("is_main", true)
	if post != nil {
		post, err = c.postService.View(post)
		if post != nil {
			if post.InfoBlocksSnapshot != nil && len(*post.InfoBlocksSnapshot) > 0 {
				err = json.Unmarshal(*post.InfoBlocksSnapshot, &blocks)
			}
		}
	}

	if err != nil {
		logger.Errorf("[blog][postController][GetHome] Error: %v", err)
	}

	c.RenderHTML(
		ctx,
		http.StatusOK,
		c.view.View(post),
		gin.H{
			"title":  "Home Page",
			"post":   post,
			"blocks": blocks,
		},
	)
}
