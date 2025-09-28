package web

import (
	"net/http"

	"github.com/axlle-com/blog/pkg/info_block/models"
	models2 "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
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
		"admin.info_block",
		gin.H{
			"title":     "Страница инфо блока",
			"templates": templates,
			"infoBlock": block,
			"settings": gin.H{
				"csrfToken": csrf.GetToken(ctx),
				"user":      user,
				"menu":      models2.NewMenu(ctx.FullPath()),
			},
		},
	)
}
