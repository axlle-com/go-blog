package web

import (
	"encoding/json"
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/dto"
	models "github.com/axlle-com/blog/pkg/blog/models"
	"github.com/gin-gonic/gin"
)

func (c *blogController) RenderPost(ctx *gin.Context, model *models.Post) {
	model, err := c.postService.View(model)
	if model == nil || err != nil {
		if err != nil {
			logger.Errorf("[blog][blogController][RenderPost] error: %v", err)
		}

		c.RenderHTML(
			ctx,
			http.StatusNotFound,
			c.view.View("error"),
			gin.H{
				"title":    "Page not found",
				"error":    "404",
				"settings": c.settings(ctx, nil),
			},
		)
		ctx.Abort()
		return
	}

	var blocks []dto.InfoBlock
	if model.InfoBlocksSnapshot != nil && len(model.InfoBlocksSnapshot) > 0 {
		if err := json.Unmarshal(model.InfoBlocksSnapshot, &blocks); err != nil {
			logger.Errorf("[blog][blogController][RenderPost] id=%v: %v", model.ID, err)
		}
	}

	c.RenderHTML(ctx,
		http.StatusOK,
		c.view.View(model.GetTemplateName()),
		gin.H{
			"settings": c.settings(ctx, model),
			"model":    model,
			"blocks":   blocks,
		},
	)
}
