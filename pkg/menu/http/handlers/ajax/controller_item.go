package ajax

import (
	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/menu/service"
	"github.com/gin-gonic/gin"
)

type ControllerMenuItem interface {
	GetMenuItems(*gin.Context)
	DeleteMenuItem(*gin.Context)
}

func NewMenuItemAjaxController(
	menuItemService *service.MenuItemService,
	menuItemCollectionService *service.MenuItemCollectionService,
	menuService *service.MenuService,
	api *api.Api,
) ControllerMenuItem {
	return &controllerItem{
		menuItemService:           menuItemService,
		menuItemCollectionService: menuItemCollectionService,
		menuService:               menuService,
		api:                       api,
	}
}

type controllerItem struct {
	*app.BaseAjax

	menuItemService           *service.MenuItemService
	menuItemCollectionService *service.MenuItemCollectionService
	menuService               *service.MenuService
	api                       *api.Api
}

func (c *controllerItem) templates(ctx *gin.Context) []contract.Template {
	templates, err := c.api.Template.GetForResources(&models.Menu{})
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}
	return templates
}
