package web

import (
	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/menu/service"
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
	api *api.Api,
) Controller {
	return &menuController{
		menuService:               menuService,
		menuCollectionService:     menuCollectionService,
		menuItemService:           menuItemService,
		menuItemCollectionService: menuItemCollectionService,
		api:                       api,
	}
}

type menuController struct {
	*app.BaseAjax

	menuService               *service.MenuService
	menuCollectionService     *service.MenuCollectionService
	menuItemService           *service.MenuItemService
	menuItemCollectionService *service.MenuItemCollectionService
	api                       *api.Api
}

func (c *menuController) templates(ctx *gin.Context) []contract.Template {
	templates, err := c.api.Template.GetForResources(&models.Menu{})
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	return templates
}
