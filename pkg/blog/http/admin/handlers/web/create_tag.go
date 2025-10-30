package web

import (
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/blog/models"
	menu "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
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

	c.RenderHTML(ctx,
		http.StatusOK,
		"admin.tag",
		gin.H{
			"title":     "Страница тэга",
			"tags":      tags,
			"templates": c.templates(ctx),
			"tag":       tag,
			"collection": gin.H{
				"infoBlocks":          c.infoBlock.GetAll(),
				"infoBlockCollection": tag.InfoBlocks,
				"relationURL":         tag.AdminURL(),
			},
			"settings": gin.H{
				"csrfToken": csrf.GetToken(ctx),
				"user":      user,
				"menu":      menu.NewMenu(ctx.FullPath()),
			},
		},
	)
}
