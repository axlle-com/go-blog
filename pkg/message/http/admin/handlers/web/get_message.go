package web

import (
	mApp "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *messageController) GetMessage(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.HTML(http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	template, err := c.messageService.GetByID(id)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	ctx.HTML(
		http.StatusOK,
		"admin.template",
		gin.H{
			"title":         "Страница шаблона",
			"user":          user,
			"templateModel": template,
			"resources":     mApp.NewResource().Resources(),
			"menu":          models.NewMenu(ctx.FullPath()),
		},
	)
}
