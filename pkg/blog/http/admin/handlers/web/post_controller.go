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

type PostController interface {
	GetPost(*gin.Context)
	GetPosts(*gin.Context)
	CreatePost(*gin.Context)
}

func NewWebPostController(
	service *service.PostService,
	services *service.PostCollectionService,
	category *service.CategoryService,
	categories *service.CategoriesService,
	tagCollectionService *service.TagCollectionService,
	template template.TemplateProvider,
	user user.UserProvider,
	gallery gallery.GalleryProvider,
	infoBlock provider.InfoBlockProvider,
) PostController {
	return &postController{
		postService:           service,
		postCollectionService: services,
		categoryService:       category,
		categoriesService:     categories,
		tagCollectionService:  tagCollectionService,
		template:              template,
		user:                  user,
		gallery:               gallery,
		infoBlock:             infoBlock,
	}
}

type postController struct {
	*app.BaseAjax

	postService           *service.PostService
	postCollectionService *service.PostCollectionService
	categoryService       *service.CategoryService
	categoriesService     *service.CategoriesService
	tagCollectionService  *service.TagCollectionService
	template              template.TemplateProvider
	user                  user.UserProvider
	gallery               gallery.GalleryProvider
	infoBlock             provider.InfoBlockProvider
}
