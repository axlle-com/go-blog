package web

import (
	app "github.com/axlle-com/blog/pkg/app/models"
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	"github.com/axlle-com/blog/pkg/post/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	GetPost(*gin.Context)
	GetPosts(*gin.Context)
	CreatePost(*gin.Context)
}

func NewWebController(
	service *service.PostService,
	services *service.PostsService,
	category *service.CategoryService,
	categories *service.CategoriesService,
	template template.TemplateProvider,
	user user.UserProvider,
	gallery gallery.GalleryProvider,
) Controller {
	return &controller{
		postService:       service,
		postsService:      services,
		categoryService:   category,
		categoriesService: categories,
		template:          template,
		user:              user,
		gallery:           gallery,
	}
}

type controller struct {
	*app.BaseAjax

	postService       *service.PostService
	postsService      *service.PostsService
	categoryService   *service.CategoryService
	categoriesService *service.CategoriesService
	template          template.TemplateProvider
	user              user.UserProvider
	gallery           gallery.GalleryProvider
}
