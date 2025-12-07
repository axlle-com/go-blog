package web

import (
	"encoding/json"
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/dto"
	"github.com/gin-gonic/gin"
)

func (c *postController) RenderHome(ctx *gin.Context) {
	var blocks []dto.InfoBlock

	post, err := c.postService.FindByParam("is_main", true)
	if post == nil || err != nil {
		c.Render404(ctx, c.view.ViewStatic("404"), nil)
		return
	}

	post, err = c.postService.View(post)
	if post == nil || err != nil {
		if err != nil {
			logger.Errorf("[blog][postController][RenderHome] Error: %v", err)
		}
		c.Render404(ctx, c.view.ViewStatic("404"), nil)
		return
	}

	if post.InfoBlocksSnapshot != nil && len(*post.InfoBlocksSnapshot) > 0 {
		err = json.Unmarshal(*post.InfoBlocksSnapshot, &blocks)
	}

	logger.Dump(c.api.MenuProvider.GetMenuString(1, ""))

	c.RenderHTML(
		ctx,
		http.StatusOK,
		c.view.View(post.GetTemplateName()),
		gin.H{
			"title":  "Home Page",
			"post":   post,
			"blocks": blocks,
		},
	)
}
