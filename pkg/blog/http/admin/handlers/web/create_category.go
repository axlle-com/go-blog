package web

import (
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/blog/models"
	models2 "github.com/axlle-com/blog/pkg/menu/models"
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

	ctx.HTML(
		http.StatusOK,
		"admin.category",
		gin.H{
			"title":      "Страница категории",
			"categories": categories,
			"templates":  c.templates(ctx),
			"category":   category,
			"settings": gin.H{
				"csrfToken": csrf.GetToken(ctx),
				"user":      user,
				"menu":      models2.NewMenu(ctx.FullPath()),
			},
		},
	)
}
