package ajax

import (
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/menu/models"
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
	postProvider contract.PostProvider,
) ControllerMenu {
	return &menuController{
		menuService:           menuService,
		menuCollectionService: menuCollectionService,
		templateProvider:      templateProvider,
		postProvider:          postProvider,
	}
}

type menuController struct {
	*app.BaseAjax

	menuService           *service.MenuService
	menuCollectionService *service.MenuCollectionService
	templateProvider      template.TemplateProvider
	postProvider          contract.PostProvider
}

func (c *menuController) templates(ctx *gin.Context) []contract.Template {
	templates, err := c.templateProvider.GetForResources(&models.Menu{})
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	return templates
}
