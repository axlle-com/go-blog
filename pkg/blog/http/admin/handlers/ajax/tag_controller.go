package ajax

import (
	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/service"
	"github.com/gin-gonic/gin"
)

type TagController interface {
	Update(*gin.Context)
	Create(*gin.Context)
	AddInfoBlock(*gin.Context)
	Delete(*gin.Context)
	DeleteImage(*gin.Context)
	Filter(*gin.Context)
}

func NewTagController(
	tagService *service.TagService,
	tagCollectionService *service.TagCollectionService,
	api *api.Api,
) TagController {
	return &tagController{
		tagService:           tagService,
		tagCollectionService: tagCollectionService,
		api:                  api,
	}
}

type tagController struct {
	*app.BaseAjax

	tagService           *service.TagService
	tagCollectionService *service.TagCollectionService
	api                  *api.Api
}

func (c *tagController) templates(ctx *gin.Context) []contract.Template {
	templates, err := c.api.Template.GetForResources(&models.PostTag{})
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	return templates
}
