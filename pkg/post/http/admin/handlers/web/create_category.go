package web

import (
	"github.com/axlle-com/blog/app/logger"
	models2 "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/post/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
)

func (c *controllerCategory) CreateCategory(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	category := &models.PostCategory{}
	categories, err := c.categoriesService.GetAll()
	if err != nil {
		logger.Error(err)
	}

	templates := c.templateProvider.GetAll()
	ctx.HTML(
		http.StatusOK,
		"admin.category",
		gin.H{
			"title":      "Страница категории",
			"categories": categories,
			"templates":  templates,
			"category":   category,
			"settings": gin.H{
				"csrfToken": csrf.GetToken(ctx),
				"user":      user,
				"menu":      models2.NewMenu(ctx.FullPath()),
			},
		},
	)
}
