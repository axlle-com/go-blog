package web

import (
	"net/http"

	app "github.com/axlle-com/blog/app/models"
	menu "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/template/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (c *messageController) CreateMessage(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	template := &models.Template{}
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
				"menu":      menu.NewMenu(ctx.FullPath(), c.GetT(ctx)),
			},
		},
	)
}
