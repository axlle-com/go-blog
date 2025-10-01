package ajax

import (
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/menu/service"
	"github.com/gin-gonic/gin"
)

type ControllerMenuItem interface {
	GetMenuItems(*gin.Context)
}

func NewMenuItemWebController(
	menuItemService *service.MenuItemService,
	menuItemCollectionService *service.MenuItemCollectionService,
) ControllerMenuItem {
	return &controller{
		menuItemService:           menuItemService,
		menuItemCollectionService: menuItemCollectionService,
	}
}

type controller struct {
	*app.BaseAjax

	menuItemService           *service.MenuItemService
	menuItemCollectionService *service.MenuItemCollectionService
}
