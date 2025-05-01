package web

import (
	mApp "github.com/axlle-com/blog/app/models"
	models2 "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/template/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
)

func (c *templateWebController) CreateTemplate(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	template := &models.Template{}
	ctx.HTML(
		http.StatusOK,
		"admin.template",
		gin.H{
			"title":         "Страница шаблона",
			"templateModel": template,
			"resources":     mApp.NewResource().Resources(),
			"settings": gin.H{
				"csrfToken": csrf.GetToken(ctx),
				"user":      user,
				"menu":      models2.NewMenu(ctx.FullPath()),
			},
		},
	)
}
