package web

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models"
	. "github.com/axlle-com/blog/pkg/info_block/models"
	models2 "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
)

func (c *infoBlockWebController) GetInfoBlocks(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}
	filter, validError := NewInfoBlockFilter().ValidateQuery(ctx)
	if validError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errors":  validError.Errors,
			"message": validError.Message,
		})
		ctx.Abort()
		return
	}
	if filter == nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка сервера"})
		return
	}
	empty := &InfoBlock{}
	paginator := models.PaginatorFromQuery(ctx.Request.URL.Query())
	paginator.SetURL(empty.AdminURL())

	templates := c.templateProvider.GetAll()
	users := c.userProvider.GetAll()

	blocksTemp, err := c.blockCollectionService.WithPaginate(paginator, filter)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}
	blocks := c.blockCollectionService.Aggregates(blocksTemp)

	ctx.HTML(http.StatusOK, "admin.info_blocks", gin.H{
		"title":      "Страница инфо блоков",
		"infoBlocks": blocks,
		"infoBlock":  empty,
		"templates":  templates,
		"users":      users,
		"paginator":  paginator,
		"filter":     filter,
		"settings": gin.H{
			"csrfToken": csrf.GetToken(ctx),
			"user":      user,
			"menu":      models2.NewMenu(ctx.FullPath()),
		},
	})
}
