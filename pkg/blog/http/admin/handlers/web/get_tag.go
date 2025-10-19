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
		ctx.HTML(http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	temp, err := c.tagService.GetByID(id)
	if err != nil {
		logger.Error(err.Error())
		ctx.HTML(http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	tag, err := c.tagService.Aggregate(temp)
	if err != nil {
		logger.Error(err.Error())
	}
	if tag == nil {
		ctx.HTML(http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	ctx.HTML(
		http.StatusOK,
		"admin.tag",
		gin.H{
			"title":     "Страница тега",
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
				"menu":      models.NewMenu(ctx.FullPath()),
			},
		},
	)
}
