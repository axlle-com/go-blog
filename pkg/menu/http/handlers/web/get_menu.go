package web

import (
	"net/http"

	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (c *menuController) GetMenu(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	id := c.GetID(ctx)
	if id == 0 {
		c.RenderHTML(ctx, http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	model, err := c.menuService.GetByID(id)
	if err != nil {
		c.RenderHTML(ctx, http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	model, err = c.menuService.Aggregate(model)
	if err != nil {
		c.RenderHTML(ctx, http.StatusInternalServerError, "admin.404", gin.H{"title": err.Error()})
		return
	}

	c.RenderHTML(ctx,
		http.StatusOK,
		"admin.menu",
		gin.H{
			"title":     c.T(ctx, "ui.name.menu"),
			"templates": c.templates(ctx),
			"model":     model,
			"settings": gin.H{
				"csrfToken": csrf.GetToken(ctx),
				"user":      user,
				"menu":      models.NewMenu(ctx.FullPath(), c.GetT(ctx)),
			},
		},
	)
}
