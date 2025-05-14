package web

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
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

	templates := c.template.GetAll()
	infoBlocks := c.infoBlock.GetAll()
	ctx.HTML(
		http.StatusOK,
		"admin.tag",
		gin.H{
			"title":     "Страница тега",
			"templates": templates,
			"tag":       tag,
			"collection": gin.H{
				"infoBlocks":         infoBlocks,
				"ifoBlockCollection": tag.InfoBlocks,
				"relationURL":        tag.AdminURL(),
			},
			"settings": gin.H{
				"csrfToken": csrf.GetToken(ctx),
				"user":      user,
				"menu":      models.NewMenu(ctx.FullPath()),
			},
		},
	)
}
