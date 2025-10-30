package web

import (
	"net/http"

	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (c *messageController) GetMessage(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		c.RenderHTML(ctx, http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	template, err := c.messageService.GetByID(id)
	if err != nil {
		c.RenderHTML(ctx, http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	c.RenderHTML(ctx,
		http.StatusOK,
		"admin.template",
		gin.H{
			"title":         "Страница шаблона",
			"templateModel": template,
			"resources":     app.NewResources().Resources(),
			"settings": gin.H{
				"csrfToken": csrf.GetToken(ctx),
				"user":      user,
				"menu":      models.NewMenu(ctx.FullPath()),
			},
		},
	)
}
