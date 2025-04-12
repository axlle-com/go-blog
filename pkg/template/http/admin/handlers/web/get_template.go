package web

import (
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *templateWebController) GetTemplate(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.HTML(http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	template, err := c.templateService.GetByID(id)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	ctx.HTML(
		http.StatusOK,
		"admin.block",
		gin.H{
			"title":    "Страница инфо блока",
			"user":     user,
			"template": template,
			"menu":     models.NewMenu(ctx.FullPath()),
		},
	)
}
