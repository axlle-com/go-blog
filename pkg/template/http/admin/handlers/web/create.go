package web

import (
	"net/http"

	app "github.com/axlle-com/blog/app/models"
	menu "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/template/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (c *templateWebController) CreateTemplate(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	template := &models.Template{}
	resources := app.NewResources()
	c.RenderHTML(ctx,
		http.StatusOK,
		"admin.template",
		gin.H{
			"title":         c.T(ctx, "ui.page.template"),
			"templateModel": template,
			"resources":     resources.Resources(),
			"themes":        resources.Themes(),
			"settings": gin.H{
				"csrfToken": csrf.GetToken(ctx),
				"user":      user,
				"menu":      menu.NewMenu(ctx.FullPath(), c.GetT(ctx)),
			},
		},
	)
}
