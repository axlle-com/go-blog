package ajax

import (
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/template/service"
	userProvider "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
)

type TemplateController interface {
	GetTemplate(ctx *gin.Context)
	UpdateTemplate(*gin.Context)
	CreateTemplate(*gin.Context)
	DeleteTemplate(*gin.Context)
	FilterTemplate(*gin.Context)
}

func NewTemplateController(
	templateService *service.TemplateService,
	templateCollectionService *service.TemplateCollectionService,
	userProvider userProvider.UserProvider,
) TemplateController {
	return &templateController{
		templateService:           templateService,
		templateCollectionService: templateCollectionService,
		userProvider:              userProvider,
	}
}

type templateController struct {
	*app.BaseAjax

	templateService           *service.TemplateService
	templateCollectionService *service.TemplateCollectionService
	userProvider              userProvider.UserProvider
}
