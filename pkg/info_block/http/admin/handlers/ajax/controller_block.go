package ajax

import (
	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/service"
	"github.com/gin-gonic/gin"
)

type InfoBlockController interface {
	GetInfoBlock(ctx *gin.Context)
	UpdateInfoBlock(*gin.Context)
	CreateInfoBlock(*gin.Context)
	DeleteInfoBlock(*gin.Context)
	DeleteBlockImage(*gin.Context)
	FilterInfoBlock(*gin.Context)
	GetInfoBlockCard(*gin.Context)
	DetachInfoBlock(*gin.Context)
}

func NewInfoBlockController(
	blockService *service.Service,
	blockCollectionService *service.CollectionService,
	api *api.Api,
) InfoBlockController {
	return &blockController{
		blockService:           blockService,
		blockCollectionService: blockCollectionService,
		api:                    api,
	}
}

type blockController struct {
	*app.BaseAjax

	blockService           *service.Service
	blockCollectionService *service.CollectionService
	api                    *api.Api
}

func (c *blockController) templates(ctx *gin.Context) []contract.Template {
	templates, err := c.api.Template.GetForResources(&models.InfoBlock{})
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	return templates
}
