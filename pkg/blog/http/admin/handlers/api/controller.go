package api

import (
	"github.com/axlle-com/blog/app/api"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/blog/service"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	GetPost(*gin.Context)
	GetPosts(*gin.Context)
	UpdatePost(*gin.Context)
	CreatePost(*gin.Context)
	DeletePost(*gin.Context)
}

func New(
	service *service.PostService,
	category *service.CategoryService,
	categories *service.CategoriesService,
	api *api.Api,
) Controller {
	return &controller{
		service:    service,
		category:   category,
		categories: categories,
		api:        api,
	}
}

type controller struct {
	*app.BaseAjax

	service    *service.PostService
	category   *service.CategoryService
	categories *service.CategoriesService
	api        *api.Api
}
