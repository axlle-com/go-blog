package web

import (
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/blog/models"
	menu "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (c *categoryController) CreateCategory(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	category := &models.PostCategory{}
	categories, err := c.categoriesService.GetAll()
	if err != nil {
		logger.Error(err)
	}

	c.RenderHTML(ctx,
		http.StatusOK,
		"admin.category",
		gin.H{
			"categories": categories,
			"templates":  c.templates(ctx),
			"category":   category,
			"collection": gin.H{
				"infoBlocks":          c.api.InfoBlock.GetAll(),
				"infoBlockCollection": category.InfoBlocks,
				"relationURL":         category.AdminURL(),
			},
			"settings": gin.H{
				"title":     c.T(ctx, "ui.name.category"),
				"csrfToken": csrf.GetToken(ctx),
				"user":      user,
				"menu":      menu.NewMenu(ctx.FullPath(), c.GetT(ctx)),
			},
		},
	)
}
