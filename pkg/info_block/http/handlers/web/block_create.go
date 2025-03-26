package web

import (
	"github.com/axlle-com/blog/pkg/info_block/models"
	models2 "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *infoBlockWebController) CreateInfoBlock(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}

	block := &models.InfoBlock{}
	templates := c.templateProvider.GetAll()
	ctx.HTML(
		http.StatusOK,
		"admin.block",
		gin.H{
			"title":     "Страница инфо блока",
			"user":      user,
			"templates": templates,
			"menu":      models2.NewMenu(ctx.FullPath()),
			"infoBlock": block,
		},
	)
}
