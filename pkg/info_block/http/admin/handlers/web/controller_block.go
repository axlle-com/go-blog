package web

import (
	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/service"
	"github.com/gin-gonic/gin"
)

type InfoBlockWebController interface {
	GetInfoBlock(*gin.Context)
	GetInfoBlocks(*gin.Context)
	CreateInfoBlock(*gin.Context)
}

func NewInfoBlockWebController(
	blockService *service.Service,
	blockCollectionService *service.CollectionService,
	api *api.Api,
) InfoBlockWebController {
	return &infoBlockWebController{
		blockService:           blockService,
		blockCollectionService: blockCollectionService,
		api:                    api,
	}
}

type infoBlockWebController struct {
	*app.BaseAjax

	blockService           *service.Service
	blockCollectionService *service.CollectionService
	api                    *api.Api
}

func (c *infoBlockWebController) templates(ctx *gin.Context) []contract.Template {
	templates, err := c.api.Template.GetForResources(&models.InfoBlock{})
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	return templates
}
