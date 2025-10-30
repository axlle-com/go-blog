package web

import (
	"net/http"

	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (c *controller) Index(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}
	c.RenderHTML(ctx, http.StatusOK,
		"admin.index",
		gin.H{
			"title": "dashboard",
			"settings": gin.H{
				"csrfToken": csrf.GetToken(ctx),
				"user":      user,
				"menu":      models.NewMenu(ctx.FullPath()),
			},
		},
	)
}
