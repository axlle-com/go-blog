package ajax

import (
	"net/http"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/gin-gonic/gin"
)

func (c *controllerItem) DeleteMenuItem(ctx *gin.Context) {
	id := c.GetID(ctx)

	if id == 0 {
		ctx.JSON(
			http.StatusOK,
			response.OK(
				response.Body{
					"status": true,
				},
				c.T(ctx, "ui.message.item_deleted"),
				nil,
			),
		)
		return
	}

	menuItem, err := c.menuItemService.GetByID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		return
	}

	menuID := menuItem.MenuID

	if err := c.menuItemService.Delete(id); err != nil {
		logger.WithRequest(ctx).Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	menu, err := c.menuService.GetByID(menuID)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	menu, err = c.menuService.Aggregate(menu)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	data := response.Body{
		"model":     menu,
		"templates": c.templates(ctx),
		"resources": app.NewResources().Resources(),
	}

	ctx.JSON(
		http.StatusOK,
		response.OK(
			response.Body{
				"view": c.RenderView("admin.menu_inner", data, ctx),
			},
			c.T(ctx, "ui.message.menu_item_deleted"),
			nil,
		),
	)
}
