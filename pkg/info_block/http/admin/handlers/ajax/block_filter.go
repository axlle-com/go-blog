package ajax

import (
	"net/http"

	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/info_block/http/admin/request"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/gin-gonic/gin"
)

func (c *blockController) FilterInfoBlock(ctx *gin.Context) {
	filter, validError := request.NewInfoBlockRequest().ValidateQuery(ctx)
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

	paginator := c.PaginatorFromQuery(ctx)
	paginator.SetURL("/admin/models-blocks")

	blocksTemp, err := c.blockCollectionService.WithPaginate(paginator, filter.ToFilter())
	if err != nil {
		logger.Error(err)
	}
	blocks := c.blockCollectionService.Aggregates(blocksTemp)

	users := c.api.User.GetAll()
	data := response.Body{
		"title":      "Страница инфо блоков",
		"infoBlocks": blocks,
		"infoBlock":  &models.InfoBlock{},
		"templates":  c.templates(ctx),
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
