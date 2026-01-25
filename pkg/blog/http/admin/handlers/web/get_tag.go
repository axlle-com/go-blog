package web

import (
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (c *tagController) GetTag(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		c.RenderHTML(ctx, http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	temp, err := c.tagService.GetByID(id)
	if err != nil {
		logger.Error(err.Error())
		c.RenderHTML(ctx, http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	tag, err := c.tagService.Aggregate(temp)
	if err != nil {
		logger.Error(err.Error())
	}
	if tag == nil {
		c.RenderHTML(ctx, http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	c.RenderHTML(ctx,
		http.StatusOK,
		"admin.tag",
		gin.H{
			"title":     c.T(ctx, "ui.name.tag"),
			"templates": c.templates(ctx),
			"tag":       tag,
			"collection": gin.H{
				"infoBlocks":          c.api.InfoBlock.GetAll(),
				"infoBlockCollection": tag.InfoBlocks,
				"relationURL":         tag.AdminURL(),
			},
			"settings": gin.H{
				"csrfToken": csrf.GetToken(ctx),
				"user":      user,
				"menu":      models.NewMenu(ctx.FullPath(), c.GetT(ctx)),
			},
		},
	)
}
