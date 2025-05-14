package web

import (
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/blog/service"
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	"github.com/axlle-com/blog/pkg/info_block/provider"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
)

type TagController interface {
	GetTag(*gin.Context)
	GetTags(*gin.Context)
	CreateTag(*gin.Context)
}

func NewWebTagController(
	tagService *service.TagService,
	tagCollectionService *service.TagCollectionService,
	template template.TemplateProvider,
	user user.UserProvider,
	gallery gallery.GalleryProvider,
	infoBlock provider.InfoBlockProvider,
) TagController {
	return &tagController{
		tagService:           tagService,
		tagCollectionService: tagCollectionService,
		template:             template,
		user:                 user,
		gallery:              gallery,
		infoBlock:            infoBlock,
	}
}

type tagController struct {
	*app.BaseAjax

	tagService           *service.TagService
	tagCollectionService *service.TagCollectionService
	template             template.TemplateProvider
	user                 user.UserProvider
	gallery              gallery.GalleryProvider
	infoBlock            provider.InfoBlockProvider
}
