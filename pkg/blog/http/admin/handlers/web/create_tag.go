package web

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/blog/models"
	menu "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
)

func (c *tagController) CreateTag(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	tag := &models.PostTag{}

	tags, err := c.tagCollectionService.GetAll()
	if err != nil {
		logger.Error(err)
	}

	templates := c.template.GetAll()
	ctx.HTML(
		http.StatusOK,
		"admin.tag",
		gin.H{
			"title":     "Страница тэга",
			"tags":      tags,
			"templates": templates,
			"tag":       tag,
			"settings": gin.H{
				"csrfToken": csrf.GetToken(ctx),
				"user":      user,
				"menu":      menu.NewMenu(ctx.FullPath()),
			},
		},
	)
}
