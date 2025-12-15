package web

import (
	"encoding/json"
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/dto"
	"github.com/gin-gonic/gin"
)

func (c *blogController) RenderHome(ctx *gin.Context) {
	var blocks []dto.InfoBlock

	post, err := c.postService.FindByParam("is_main", true)
	if post == nil || err != nil {
		c.Render404(ctx, c.view.ViewStatic("error"), nil)
		return
	}

	post, err = c.postService.View(post)
	if post == nil || err != nil {
		if err != nil {
			logger.Errorf("[blog][blogController][RenderHome] Error: %v", err)
		}
		c.Render404(ctx, c.view.ViewStatic("error"), nil)
		return
	}

	if len(post.InfoBlocksSnapshot) > 0 {
		err = json.Unmarshal(post.InfoBlocksSnapshot, &blocks)
	}

	c.RenderHTML(
		ctx,
		http.StatusOK,
		c.view.View(post.GetTemplateName()),
		gin.H{
			"settings": c.settings(ctx),
			"title":    "Home Page",
			"post":     post,
			"blocks":   blocks,
		},
	)
}
