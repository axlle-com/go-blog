package web

import (
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
)

func (c *controller) Index(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}
	ctx.HTML(
		http.StatusOK,
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
