package ajax

import (
	"net/http"

	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/menu/http/request"
	"github.com/gin-gonic/gin"
)

func (c *controllerItem) GetMenuItems(ctx *gin.Context) {
	user := c.GetUser(ctx)
	if user == nil {
		return
	}
	filter, validError := request.NewMenuItemFilter().ValidateQuery(ctx)
	if validError != nil {
		ctx.JSON(http.StatusBadRequest,
			response.Fail(http.StatusBadRequest, validError.Message, validError.Errors),
		)
		ctx.Abort()
		return
	}

	paginator := c.PaginatorFromQuery(ctx)
	paginator.SetURL("/admin/menu/items")

	items, err := c.menuItemCollectionService.Filter(paginator, filter.ToFilter())
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	ctx.JSON(
		http.StatusOK,
		response.OK(
			response.Body{
				"items":     items,
				"paginator": paginator,
			},
			"",
			paginator,
		),
	)
}
