package web

import (
	app "github.com/axlle-com/blog/pkg/app/models"
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	"github.com/axlle-com/blog/pkg/post/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
)

type PostController interface {
	GetHome(*gin.Context)
	GetPost(*gin.Context)
}

func NewFrontWebController(
	service *service.PostService,
	services *service.PostsService,
	category *service.CategoryService,
	categories *service.CategoriesService,
	template template.TemplateProvider,
	user user.UserProvider,
	gallery gallery.GalleryProvider,
) PostController {
	return &postController{
		postService:       service,
		postsService:      services,
		categoryService:   category,
		categoriesService: categories,
		template:          template,
		user:              user,
		gallery:           gallery,
	}
}

type postController struct {
	*app.BaseAjax

	postService       *service.PostService
	postsService      *service.PostsService
	categoryService   *service.CategoryService
	categoriesService *service.CategoriesService
	template          template.TemplateProvider
	user              user.UserProvider
	gallery           gallery.GalleryProvider
}
