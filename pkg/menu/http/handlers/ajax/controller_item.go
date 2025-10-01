package ajax

import (
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/menu/service"
	"github.com/gin-gonic/gin"
)

type ControllerMenuItem interface {
	GetMenuItems(*gin.Context)
}

func NewMenuItemAjaxController(
	menuItemService *service.MenuItemService,
	menuItemCollectionService *service.MenuItemCollectionService,
) ControllerMenuItem {
	return &controllerItem{
		menuItemService:           menuItemService,
		menuItemCollectionService: menuItemCollectionService,
	}
}

type controllerItem struct {
	*app.BaseAjax

	menuItemService           *service.MenuItemService
	menuItemCollectionService *service.MenuItemCollectionService
}
