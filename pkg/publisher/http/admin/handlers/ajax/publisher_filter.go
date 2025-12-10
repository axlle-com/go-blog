package ajax

import (
	"net/http"

	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/publisher/models"
	"github.com/gin-gonic/gin"
)

func (c *controller) Filter(ctx *gin.Context) {
	paginator := c.PaginatorFromQuery(ctx)
	paginator.SetURL("/admin/publishers")

	filter, validError := models.NewPublisherFilter().ValidateQuery(ctx)
	if validError != nil {
		ctx.JSON(
			http.StatusBadRequest,
			response.Fail(http.StatusBadRequest, validError.Message, validError.Errors),
		)
		ctx.Abort()
		return
	}

	publishers, err := c.collectionService.WithPaginate(paginator, filter)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	ctx.JSON(
		http.StatusOK,
		response.OK(
			response.Body{
				"items":     publishers,
				"paginator": paginator,
			},
			"",
			paginator,
		),
	)
}
