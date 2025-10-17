package ajax

import (
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/menu/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	"github.com/gin-gonic/gin"
)

type ControllerMenu interface {
	Create(*gin.Context)
	Update(*gin.Context)
}

func NewMenuAjaxController(
	menuService *service.MenuService,
	menuCollectionService *service.MenuCollectionService,
	templateProvider template.TemplateProvider,
	postProvider contracts.PostProvider,
) ControllerMenu {
	return &controllerMenu{
		menuService:           menuService,
		menuCollectionService: menuCollectionService,
		templateProvider:      templateProvider,
		postProvider:          postProvider,
	}
}

type controllerMenu struct {
	*app.BaseAjax

	menuService           *service.MenuService
	menuCollectionService *service.MenuCollectionService
	templateProvider      template.TemplateProvider
	postProvider          contracts.PostProvider
}
