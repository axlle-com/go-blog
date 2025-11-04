package web

import (
	"github.com/axlle-com/blog/app/api"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/template/service"
	"github.com/gin-gonic/gin"
)

type TemplateWebController interface {
	GetTemplate(*gin.Context)
	GetTemplates(*gin.Context)
	CreateTemplate(*gin.Context)
}

func NewTemplateWebController(
	templateService *service.TemplateService,
	templateCollectionService *service.TemplateCollectionService,
	api *api.Api,
) TemplateWebController {
	return &templateWebController{
		templateService:           templateService,
		templateCollectionService: templateCollectionService,
		api:                       api,
	}
}

type templateWebController struct {
	*app.BaseAjax

	templateService           *service.TemplateService
	templateCollectionService *service.TemplateCollectionService
	api                       *api.Api
}
