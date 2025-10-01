package web

import (
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (c *controller) CreateMenu(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	model := &models.Menu{}
	templates, err := c.templateProvider.GetForResources(model)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}
	ctx.HTML(
		http.StatusOK,
		"admin.menu",
		gin.H{
			"title":     "Страница меню",
			"templates": templates,
			"model":     model,
			"settings": gin.H{
				"csrfToken": csrf.GetToken(ctx),
				"user":      user,
				"menu":      models.NewMenu(ctx.FullPath()),
			},
		},
	)
}
