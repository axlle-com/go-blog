package api

import (
	app "github.com/axlle-com/blog/pkg/app/models"
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	"github.com/axlle-com/blog/pkg/post/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	GetPost(c *gin.Context)
	GetPosts(c *gin.Context)
	UpdatePost(*gin.Context)
	CreatePost(*gin.Context)
	DeletePost(*gin.Context)
}

func New(
	service *service.PostService,
	category *service.CategoryService,
	categories *service.CategoriesService,
	template template.TemplateProvider,
	user user.UserProvider,
	gallery gallery.GalleryProvider,
) Controller {
	return &controller{
		service:    service,
		category:   category,
		categories: categories,
		template:   template,
		user:       user,
		gallery:    gallery,
	}
}

type controller struct {
	*app.BaseAjax

	service    *service.PostService
	category   *service.CategoryService
	categories *service.CategoriesService
	template   template.TemplateProvider
	user       user.UserProvider
	gallery    gallery.GalleryProvider
}
