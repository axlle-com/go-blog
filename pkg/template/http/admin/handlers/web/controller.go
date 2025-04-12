package web

import (
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/template/service"
	userProvider "github.com/axlle-com/blog/pkg/user/provider"
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
	userProvider userProvider.UserProvider,
) TemplateWebController {
	return &templateWebController{
		templateService:           templateService,
		templateCollectionService: templateCollectionService,
		userProvider:              userProvider,
	}
}

type templateWebController struct {
	*app.BaseAjax

	templateService           *service.TemplateService
	templateCollectionService *service.TemplateCollectionService
	userProvider              userProvider.UserProvider
}
