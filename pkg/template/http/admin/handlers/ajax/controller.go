package ajax

import (
	"github.com/axlle-com/blog/app/api"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/template/service"
	"github.com/gin-gonic/gin"
)

type TemplateController interface {
	GetTemplate(ctx *gin.Context)
	GetResourceTemplate(ctx *gin.Context)
	UpdateTemplate(*gin.Context)
	CreateTemplate(*gin.Context)
	DeleteTemplate(*gin.Context)
	FilterTemplate(*gin.Context)
}

func NewTemplateController(
	templateService *service.Service,
	templateCollectionService *service.CollectionService,
	api *api.Api,
) TemplateController {
	return &templateController{
		templateService:           templateService,
		templateCollectionService: templateCollectionService,
		api:                       api,
	}
}

type templateController struct {
	*app.BaseAjax

	templateService           *service.Service
	templateCollectionService *service.CollectionService
	api                       *api.Api
}
