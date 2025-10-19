package web

import (
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (c *categoryController) GetCategory(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.HTML(http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	category, err := c.categoryService.GetAggregateByID(id)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	category.Galleries = c.galleryProvider.GetForResource(category)

	categories, err := c.categoriesService.GetAllForParent(category)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	ctx.HTML(
		http.StatusOK,
		"admin.category",
		gin.H{
			"title":      "Страница категории",
			"categories": categories,
			"templates":  c.templates(ctx),
			"category":   category,
			"collection": gin.H{
				"infoBlocks":          c.infoBlockProvider.GetAll(),
				"infoBlockCollection": category.InfoBlocks,
				"relationURL":         category.AdminURL(),
			},
			"settings": gin.H{
				"csrfToken": csrf.GetToken(ctx),
				"user":      user,
				"menu":      models.NewMenu(ctx.FullPath()),
			},
		},
	)
}
