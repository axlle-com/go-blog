package web

import (
	mApp "github.com/axlle-com/blog/app/models"
	models2 "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/template/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *messageController) CreateMessage(ctx *gin.Context) {
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
			"user":          user,
			"templateModel": template,
			"resources":     mApp.NewResource().Resources(),
			"menu":          models2.NewMenu(ctx.FullPath()),
		},
	)
}
