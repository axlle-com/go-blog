package web

import (
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	"github.com/gin-gonic/gin"
)

func (c *blogController) RenderHome(ctx *gin.Context) {
	model, err := c.postService.FindByParam("is_main", true)
	if model == nil || err != nil {
		c.RenderHTML(
			ctx,
			http.StatusNotFound,
			c.view.ViewStatic("error"),
			gin.H{
				"title":    "Page not found",
				"error":    "404",
				"settings": c.settings(ctx, nil),
			},
		)
		ctx.Abort()
		return
	}

	model, err = c.postService.View(model)
	if model == nil || err != nil {
		if err != nil {
			logger.Errorf("[blog][blogController][RenderHome] Error: %v", err)
		}
		c.RenderHTML(
			ctx,
			http.StatusNotFound,
			c.view.ViewStatic("error"),
			gin.H{
				"title":    "Page not found",
				"error":    "404",
				"settings": c.settings(ctx, nil),
			},
		)
		ctx.Abort()
		return
	}

	c.RenderHTML(
		ctx,
		http.StatusOK,
		c.view.ViewResource(model),
		gin.H{
			"settings": c.settings(ctx, model),
			"model":    model,
			"blocks":   model.InfoBlocks,
		},
	)
}
