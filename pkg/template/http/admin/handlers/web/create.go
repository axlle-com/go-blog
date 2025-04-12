package web

import (
	models2 "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/template/models"
	"github.com/gin-gonic/gin"
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
		"admin.block",
		gin.H{
			"title":    "Страница инфо блока",
			"user":     user,
			"template": template,
			"menu":     models2.NewMenu(ctx.FullPath()),
		},
	)
}
