package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/menu/models"
)

func (c *menuController) GetMenu(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.HTML(http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	model, err := c.menuService.GetByID(id)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	model, err = c.menuService.Aggregate(model)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "admin.404", gin.H{"title": err.Error()})
		return
	}

	publishers, err := c.postProvider.GetPublishers()
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	ctx.HTML(
		http.StatusOK,
		"admin.menu",
		gin.H{
			"title":      "Страница меню",
			"templates":  c.templates(ctx),
			"model":      model,
			"publishers": publishers,
			"settings": gin.H{
				"csrfToken": csrf.GetToken(ctx),
				"user":      user,
				"menu":      models.NewMenu(ctx.FullPath()),
			},
		},
	)
}
