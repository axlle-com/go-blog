package web

import (
	"net/http"

	"github.com/axlle-com/blog/pkg/info_block/models"
	menuModels "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (c *infoBlockWebController) GetInfoBlock(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		c.RenderHTML(ctx, http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	block, err := c.blockService.FindByID(id)
	if err != nil {
		c.RenderHTML(ctx, http.StatusNotFound, "admin.404", gin.H{"title": "404 Not Found"})
		return
	}

	block.Galleries = c.api.Gallery.GetForResourceUUID(block.UUID.String())

	var infoBlocks []*models.InfoBlock
	if block.ID != 0 {
		var err2 error
		infoBlocks, err2 = c.blockCollectionService.GetAllForParent(block)
		if err2 != nil {
			// Если ошибка, получаем все инфоблоки
			infoBlocks, _ = c.blockCollectionService.GetAll()
		}
	} else {
		infoBlocks, _ = c.blockCollectionService.GetAll()
	}

	c.RenderHTML(ctx,
		http.StatusOK,
		"admin.info_block",
		gin.H{
			"title":      c.T(ctx, "ui.page.info_block"),
			"templates":  c.templates(ctx),
			"infoBlocks": infoBlocks,
			"infoBlock":  block,
			"settings": gin.H{
				"csrfToken": csrf.GetToken(ctx),
				"user":      user,
				"menu":      menuModels.NewMenu(ctx.FullPath(), c.GetT(ctx)),
			},
		},
	)
}
