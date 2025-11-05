package web

import (
	"github.com/axlle-com/blog/app/api"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/blog/service"
	"github.com/gin-gonic/gin"
)

type PostController interface {
	RenderHome(*gin.Context)
	FindByAlias(ctx *gin.Context)
}

func NewFrontWebController(
	view contract.View,
	service *service.PostService,
	services *service.PostCollectionService,
	category *service.CategoryService,
	categories *service.CategoriesService,
	api *api.Api,
) PostController {
	return &postController{
		view:                  view,
		postService:           service,
		postCollectionService: services,
		categoryService:       category,
		categoriesService:     categories,
		api:                   api,
	}
}

type postController struct {
	*app.BaseAjax

	view                  contract.View
	postService           *service.PostService
	postCollectionService *service.PostCollectionService
	categoryService       *service.CategoryService
	categoriesService     *service.CategoriesService
	api                   *api.Api
}
