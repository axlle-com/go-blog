package web

import (
	"encoding/json"
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/dto"
	models "github.com/axlle-com/blog/pkg/blog/models"
	"github.com/gin-gonic/gin"
)

func (c *blogController) RenderCategory(ctx *gin.Context, category *models.PostCategory) {
	filter, validError := models.NewPostFilter().ValidateQuery(ctx)
	if validError != nil {
		logger.WithRequest(ctx).Error(validError)
		filter = models.NewPostFilter()
	}
	if filter == nil {
		filter = models.NewPostFilter()
	}

	paginator := c.PaginatorFromQuery(ctx)
	paginator.SetURL(category.GetURL())

	category, err := c.categoryService.View(category, paginator, filter)
	if category == nil || err != nil {
		if err != nil {
			logger.Errorf("[blog][blogController][RenderCategory] error: %v", err)
		}

		c.Render404(ctx, c.view.View("error"), nil)
		return
	}

	var blocks []dto.InfoBlock
	if len(category.InfoBlocksSnapshot) > 0 {
		if err := json.Unmarshal(category.InfoBlocksSnapshot, &blocks); err != nil {
			logger.Errorf("[blog][blogController][RenderCategory] id=%v: %v", category.ID, err)
		}
	}

	c.RenderHTML(ctx,
		http.StatusOK,
		c.view.View(category.GetTemplateName()),
		gin.H{
			"settings":  c.settings(ctx),
			"title":     category.Title,
			"blocks":    blocks,
			"category":  category,
			"paginator": paginator,
			"filter":    filter,
		},
	)
}
