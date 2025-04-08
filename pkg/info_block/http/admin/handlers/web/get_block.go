package web

import (
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *infoBlockWebController) GetInfoBlock(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.HTML(http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	block, err := c.blockService.GetByID(id)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	block.Galleries = c.galleryProvider.GetForResource(block)
	templates := c.templateProvider.GetAll()
	ctx.HTML(
		http.StatusOK,
		"admin.block",
		gin.H{
			"title":     "Страница инфо блока",
			"user":      user,
			"templates": templates,
			"menu":      models.NewMenu(ctx.FullPath()),
			"infoBlock": block,
		},
	)
}
