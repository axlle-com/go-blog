package web

import (
	"github.com/axlle-com/blog/pkg/app/logger"
	models2 "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/post/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controllerCategory) CreateCategory(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	category := &models.PostCategory{}
	categories, err := c.categoryRepo.GetAll()
	if err != nil {
		logger.Error(err)
	}

	templates := c.template.GetAll()
	ctx.HTML(
		http.StatusOK,
		"admin.category",
		gin.H{
			"title":      "Страница категории",
			"user":       user,
			"categories": categories,
			"templates":  templates,
			"menu":       models2.NewMenu(ctx.FullPath()),
			"category":   category,
		},
	)
}
