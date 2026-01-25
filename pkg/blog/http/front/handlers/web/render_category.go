package web

import (
	"encoding/json"
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/dto"
	models "github.com/axlle-com/blog/pkg/blog/models"
	"github.com/gin-gonic/gin"
)

func (c *blogController) RenderCategory(ctx *gin.Context, model *models.PostCategory) {
	filter, validError := models.NewPostFilter().ValidateQuery(ctx)
	if validError != nil {
		logger.WithRequest(ctx).Error(validError)
		filter = models.NewPostFilter()
	}
	if filter == nil {
		filter = models.NewPostFilter()
	}

	paginator := c.PaginatorFromQuery(ctx)
	paginator.SetURL(model.GetURL())

	model, err := c.categoryService.View(model, paginator, filter)
	if model == nil || err != nil {
		if err != nil {
			logger.Errorf("[blog][blogController][RenderCategory] error: %v", err)
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

	var blocks []dto.InfoBlock
	if len(model.InfoBlocksSnapshot) > 0 {
		if err := json.Unmarshal(model.InfoBlocksSnapshot, &blocks); err != nil {
			logger.Errorf("[blog][blogController][RenderCategory] id=%v: %v", model.ID, err)
		}
	}

	c.RenderHTML(
		ctx,
		http.StatusOK,
		c.view.ViewResource(model),
		gin.H{
			"settings":  c.settings(ctx, model),
			"model":     model,
			"blocks":    blocks,
			"paginator": paginator,
			"filter":    filter,
		},
	)
}
