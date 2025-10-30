package ajax

import (
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	appPovider "github.com/axlle-com/blog/app/models/provider"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/service"
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
	infoBlock appPovider.InfoBlockProvider,
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
	infoBlock            appPovider.InfoBlockProvider
}

func (c *tagController) templates(ctx *gin.Context) []contract.Template {
	templates, err := c.templateProvider.GetForResources(&models.PostTag{})
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	return templates
}
