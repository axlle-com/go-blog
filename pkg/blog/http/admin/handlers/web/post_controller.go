package web

import (
	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/service"
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
	api *api.Api,
) PostController {
	return &postController{
		postService:           service,
		postCollectionService: services,
		categoryService:       category,
		categoriesService:     categories,
		tagCollectionService:  tagCollectionService,
		api:                   api,
	}
}

type postController struct {
	*app.BaseAjax

	postService           *service.PostService
	postCollectionService *service.PostCollectionService
	categoryService       *service.CategoryService
	categoriesService     *service.CategoriesService
	tagCollectionService  *service.TagCollectionService
	api                   *api.Api
}

func (c *postController) templates(ctx *gin.Context) []contract.Template {
	templates, err := c.api.Template.GetForResources(&models.Post{})
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	return templates
}
