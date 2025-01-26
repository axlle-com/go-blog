package web

import (
	app "github.com/axlle-com/blog/pkg/app/models"
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	"github.com/axlle-com/blog/pkg/post/repository"
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
	service *service.Service,
	post repository.PostRepository,
	category repository.CategoryRepository,
	template template.TemplateProvider,
	user user.UserProvider,
	gallery gallery.GalleryProvider,
) Controller {
	return &controller{
		service:  service,
		post:     post,
		category: category,
		template: template,
		user:     user,
		gallery:  gallery,
	}
}

type controller struct {
	*app.BaseAjax

	service  *service.Service
	post     repository.PostRepository
	category repository.CategoryRepository
	template template.TemplateProvider
	user     user.UserProvider
	gallery  gallery.GalleryProvider
}
