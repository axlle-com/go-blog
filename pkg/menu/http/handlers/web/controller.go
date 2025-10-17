package web

import (
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/menu/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	GetMenu(*gin.Context)
	GetMenus(*gin.Context)
	CreateMenu(*gin.Context)
}

func NewMenuWebController(
	menuService *service.MenuService,
	menuCollectionService *service.MenuCollectionService,
	menuItemService *service.MenuItemService,
	menuItemCollectionService *service.MenuItemCollectionService,
	templateProvider template.TemplateProvider,
	postProvider contracts.PostProvider,
) Controller {
	return &controller{
		menuService:               menuService,
		menuCollectionService:     menuCollectionService,
		menuItemService:           menuItemService,
		menuItemCollectionService: menuItemCollectionService,
		templateProvider:          templateProvider,
		postProvider:              postProvider,
	}
}

type controller struct {
	*app.BaseAjax

	menuService               *service.MenuService
	menuCollectionService     *service.MenuCollectionService
	menuItemService           *service.MenuItemService
	menuItemCollectionService *service.MenuItemCollectionService
	templateProvider          template.TemplateProvider
	postProvider              contracts.PostProvider
}
