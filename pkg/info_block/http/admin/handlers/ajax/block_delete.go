package ajax

import (
	"net/http"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/info_block/http/admin/request"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/gin-gonic/gin"
)

func (c *blockController) DeleteInfoBlock(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	block, err := c.blockService.FindByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	if err := c.blockService.Delete(block); err != nil {
		logger.WithRequest(ctx).Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

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
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": c.T(ctx, "ui.error.server_error")})
		return
	}

	paginator := c.PaginatorFromQuery(ctx)
	paginator.SetURL("/admin/models-blocks")

	users := c.api.User.GetAll()

	blocksTemp, err := c.blockCollectionService.WithPaginate(paginator, filter.ToFilter())
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}
	blocks := c.blockCollectionService.Aggregates(blocksTemp)

	data := response.Body{
		"title":      c.T(ctx, "ui.page.info_blocks"),
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
				"view": c.RenderView("admin.info_blocks_inner", data, ctx),
			},
			c.T(ctx, "ui.success.record_deleted"),
			paginator,
		),
	)
}
