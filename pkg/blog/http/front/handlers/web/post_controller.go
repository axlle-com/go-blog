package web

import (
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/blog/service"
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
)

type PostController interface {
	GetHome(*gin.Context)
	GetPost(*gin.Context)
}

func NewFrontWebController(
	view contracts.View,
	service *service.PostService,
	services *service.PostCollectionService,
	category *service.CategoryService,
	categories *service.CategoriesService,
	template template.TemplateProvider,
	user user.UserProvider,
	gallery gallery.GalleryProvider,
) PostController {
	return &postController{
		view:                  view,
		postService:           service,
		postCollectionService: services,
		categoryService:       category,
		categoriesService:     categories,
		template:              template,
		user:                  user,
		gallery:               gallery,
	}
}

type postController struct {
	*app.BaseAjax

	view                  contracts.View
	postService           *service.PostService
	postCollectionService *service.PostCollectionService
	categoryService       *service.CategoryService
	categoriesService     *service.CategoriesService
	template              template.TemplateProvider
	user                  user.UserProvider
	gallery               gallery.GalleryProvider
}
