package ajax

import (
	app "github.com/axlle-com/blog/pkg/app/models"
	"github.com/axlle-com/blog/pkg/post/repository"
	"github.com/axlle-com/blog/pkg/post/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	UpdatePost(*gin.Context)
	CreatePost(*gin.Context)
	DeletePost(*gin.Context)
	DeletePostImage(*gin.Context)
	FilterPosts(*gin.Context)
}

func New(
	service *service.Service,
	post repository.PostRepository,
	category repository.CategoryRepository,
	template template.TemplateProvider,
	user user.UserProvider,
) Controller {
	return &controller{
		service:  service,
		post:     post,
		category: category,
		template: template,
		user:     user,
	}
}

type controller struct {
	*app.BaseAjax

	service  *service.Service
	post     repository.PostRepository
	category repository.CategoryRepository
	template template.TemplateProvider
	user     user.UserProvider
}
