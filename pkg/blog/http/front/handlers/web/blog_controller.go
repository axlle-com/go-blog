package web

import (
	"html/template"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/blog/service"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

type BlogController interface {
	RenderHome(*gin.Context)
	RenderByURL(ctx *gin.Context)
}

func NewFrontWebController(
	view contract.View,
	service *service.PostService,
	services *service.PostCollectionService,
	category *service.CategoryService,
	categories *service.CategoriesService,
	api *api.Api,
) BlogController {
	return &blogController{
		view:                  view,
		postService:           service,
		postCollectionService: services,
		categoryService:       category,
		categoriesService:     categories,
		api:                   api,
	}
}

type blogController struct {
	*app.BaseAjax

	view                  contract.View
	postService           *service.PostService
	postCollectionService *service.PostCollectionService
	categoryService       *service.CategoryService
	categoriesService     *service.CategoriesService
	api                   *api.Api
}

func (c *blogController) settings(ctx *gin.Context) map[string]any {
	menu, err := c.api.Menu.GetMenuString(1, ctx.Request.URL.Path)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	return gin.H{
		"menu":      template.HTML(menu),
		"user":      c.GetAdmin(ctx),
		"csrfToken": csrf.GetToken(ctx),
	}
}
