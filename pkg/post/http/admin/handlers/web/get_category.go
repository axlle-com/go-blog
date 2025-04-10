package web

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controllerCategory) GetCategory(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.HTML(http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	category, err := c.categoryService.GetByID(id)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	category.Galleries = c.galleryProvider.GetForResource(category)

	categories, err := c.categoriesService.GetAllForParent(category)
	if err != nil {
		logger.Error(err)
	}

	templates := c.templateProvider.GetAll()
	ctx.HTML(
		http.StatusOK,
		"admin.category",
		gin.H{
			"title":        "Страница категории",
			"userProvider": user,
			"categories":   categories,
			"templates":    templates,
			"menu":         models.NewMenu(ctx.FullPath()),
			"category":     category,
		},
	)
}
