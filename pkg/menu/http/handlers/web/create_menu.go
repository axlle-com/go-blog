package web

import (
	"net/http"

	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (c *menuController) CreateMenu(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	c.RenderHTML(ctx,
		http.StatusOK,
		"admin.menu",
		gin.H{
			"title":     c.T(ctx, "ui.page.menu"),
			"templates": c.templates(ctx),
			"model":     &models.Menu{},
			"settings": gin.H{
				"csrfToken": csrf.GetToken(ctx),
				"user":      user,
				"menu":      models.NewMenu(ctx.FullPath(), c.GetT(ctx)),
			},
		},
	)
}
