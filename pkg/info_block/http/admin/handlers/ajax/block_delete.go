package ajax

import (
	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models"
	. "github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *blockController) DeleteInfoBlock(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	block, err := c.blockService.GetByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	if err := c.blockService.Delete(block); err != nil {
		logger.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	filter, validError := NewInfoBlockFilter().ValidateQuery(ctx)
	if validError != nil {
		ctx.JSON(
			http.StatusBadRequest,
			response.Fail(http.StatusBadRequest, validError.Message, validError.Errors),
		)
		ctx.Abort()
		return
	}

	if filter == nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка сервера"})
		return
	}

	paginator := models.PaginatorFromQuery(ctx.Request.URL.Query())
	paginator.SetURL("/admin/info-blocks")

	users := c.userProvider.GetAll()
	templates := c.templateProvider.GetAll()

	blocksTemp, err := c.blockCollectionService.WithPaginate(paginator, filter)
	if err != nil {
		logger.Error(err)
	}
	blocks := c.blockCollectionService.Aggregates(blocksTemp)

	data := response.Body{
		"title":      "Страница инфо блоков",
		"infoBlocks": blocks,
		"infoBlock":  &InfoBlock{},
		"templates":  templates,
		"users":      users,
		"paginator":  paginator,
		"filter":     filter,
	}

	ctx.JSON(
		http.StatusOK,
		response.OK(
			response.Body{
				"view": c.RenderView("admin.blocks_inner", data, ctx),
			},
			"Запись удалена",
			paginator,
		),
	)
}
