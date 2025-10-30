package web

import (
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	appPovider "github.com/axlle-com/blog/app/models/provider"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
)

type InfoBlockWebController interface {
	GetInfoBlock(*gin.Context)
	GetInfoBlocks(*gin.Context)
	CreateInfoBlock(*gin.Context)
}

func NewInfoBlockWebController(
	blockService *service.InfoBlockService,
	blockCollectionService *service.InfoBlockCollectionService,
	template template.TemplateProvider,
	user user.UserProvider,
	galleryProvider appPovider.GalleryProvider,
) InfoBlockWebController {
	return &infoBlockWebController{
		blockService:           blockService,
		blockCollectionService: blockCollectionService,
		templateProvider:       template,
		userProvider:           user,
		galleryProvider:        galleryProvider,
	}
}

type infoBlockWebController struct {
	*app.BaseAjax

	blockService           *service.InfoBlockService
	blockCollectionService *service.InfoBlockCollectionService
	templateProvider       template.TemplateProvider
	userProvider           user.UserProvider
	galleryProvider        appPovider.GalleryProvider
}

func (c *infoBlockWebController) templates(ctx *gin.Context) []contract.Template {
	templates, err := c.templateProvider.GetForResources(&models.InfoBlock{})
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	return templates
}
