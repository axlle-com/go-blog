package ajax

import (
	"net/http"

	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	"github.com/gin-gonic/gin"
)

func (c *controller) Filter(ctx *gin.Context) {
	paginator := c.PaginatorFromQuery(ctx)
	paginator.SetURL("/admin/publishers")

	publishers, err := c.collectionService.WithPaginate(paginator)
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
