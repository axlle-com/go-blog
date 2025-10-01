package ajax

import (
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/menu/service"
	"github.com/gin-gonic/gin"
)

type ControllerMenu interface {
	Create(*gin.Context)
	Update(*gin.Context)
}

func NewMenuAjaxController(
	menuService *service.MenuService,
	menuCollectionService *service.MenuCollectionService,
) ControllerMenu {
	return &controllerMenu{
		menuService:           menuService,
		menuCollectionService: menuCollectionService,
	}
}

type controllerMenu struct {
	*app.BaseAjax

	menuService           *service.MenuService
	menuCollectionService *service.MenuCollectionService
}
