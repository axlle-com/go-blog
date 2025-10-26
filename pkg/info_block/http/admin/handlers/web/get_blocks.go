package web

import (
	"net/http"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/info_block/http/admin/request"
	"github.com/axlle-com/blog/pkg/info_block/models"
	menu "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (c *infoBlockWebController) GetInfoBlocks(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}
	filter, validError := request.NewInfoBlockRequest().ValidateQuery(ctx)
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
	empty := &models.InfoBlock{}
	paginator := c.PaginatorFromQuery(ctx)
	paginator.SetURL(empty.AdminURL())

	users := c.userProvider.GetAll()

	blocksTemp, err := c.blockCollectionService.WithPaginate(paginator, filter.ToFilter())
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}
	blocks := c.blockCollectionService.Aggregates(blocksTemp)

	ctx.HTML(http.StatusOK, "admin.info_blocks", gin.H{
		"title":      "Страница инфо блоков",
		"infoBlocks": blocks,
		"infoBlock":  empty,
		"templates":  c.templates(ctx),
		"users":      users,
		"paginator":  paginator,
		"filter":     filter,
		"settings": gin.H{
			"csrfToken": csrf.GetToken(ctx),
			"user":      user,
			"menu":      menu.NewMenu(ctx.FullPath()),
		},
	})
}
