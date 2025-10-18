package ajax

import (
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/service"
	infoBlock "github.com/axlle-com/blog/pkg/info_block/provider"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
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
	template template.TemplateProvider,
	user user.UserProvider,
	infoBlock infoBlock.InfoBlockProvider,
) TagController {
	return &tagController{
		tagService:           tagService,
		tagCollectionService: tagCollectionService,
		templateProvider:     template,
		user:                 user,
		infoBlock:            infoBlock,
	}
}

type tagController struct {
	*app.BaseAjax

	tagService           *service.TagService
	tagCollectionService *service.TagCollectionService
	templateProvider     template.TemplateProvider
	user                 user.UserProvider
	infoBlock            infoBlock.InfoBlockProvider
}

func (c *tagController) templates(ctx *gin.Context) []contracts.Template {
	templates, err := c.templateProvider.GetForResources(&models.PostTag{})
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	return templates
}
