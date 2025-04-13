package ajax

import (
	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models"
	. "github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *blockController) FilterInfoBlock(ctx *gin.Context) {
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
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Message("Ошибка сервера"))
		return
	}

	paginator := models.PaginatorFromQuery(ctx.Request.URL.Query())
	paginator.SetURL("/admin/info-blocks")

	blocksTemp, err := c.blockCollectionService.WithPaginate(paginator, filter)
	if err != nil {
		logger.Error(err)
	}
	blocks := c.blockCollectionService.Aggregates(blocksTemp)

	templates := c.templateProvider.GetAll()
	users := c.userProvider.GetAll()
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
				"infoBlocks": blocks,
				"paginator":  paginator,
				"url":        filter.GetURL(),
				"view":       c.RenderView("admin.info_blocks_inner", data, ctx),
			},
			"",
			paginator,
		),
	)
}
