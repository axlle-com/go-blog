package web

import (
	"encoding/json"
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/dto"
	models "github.com/axlle-com/blog/pkg/blog/models"
	"github.com/gin-gonic/gin"
)

func (c *postController) RenderPost(ctx *gin.Context, post *models.Post) {
	post, err := c.postService.View(post)
	if post == nil || err != nil {
		if err != nil {
			logger.Errorf("[blog][postController][RenderPost] error: %v", err)
		}

		c.Render404(ctx, c.view.ViewStatic("404"), nil)
		return
	}

	var blocks []dto.InfoBlock
	if post.InfoBlocksSnapshot != nil && len(*post.InfoBlocksSnapshot) > 0 {
		if err := json.Unmarshal(*post.InfoBlocksSnapshot, &blocks); err != nil {
			logger.Errorf("[blog][postController][RenderPost] id=%v: %v", post.ID, err)
		}
	}

	c.RenderHTML(ctx,
		http.StatusOK,
		c.view.View(post.GetTemplateName()),
		gin.H{
			"title":  "Home Page",
			"blocks": blocks,
		},
	)
}
