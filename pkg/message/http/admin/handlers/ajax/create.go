package web

import (
	"net/http"

	app "github.com/axlle-com/blog/app/models"
	menu "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/template/models"
	"github.com/gin-gonic/gin"
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
			"user":          user,
			"templateModel": template,
			"resources":     app.NewResources().Resources(),
			"menu":          menu.NewMenu(ctx.FullPath(), c.BuildT(ctx)),
		},
	)
}
