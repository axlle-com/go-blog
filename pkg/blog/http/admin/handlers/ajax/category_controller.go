package ajax

import (
	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/service"
	"github.com/gin-gonic/gin"
)

type CategoryController interface {
	UpdateCategory(*gin.Context)
	CreateCategory(*gin.Context)
	DeleteCategory(*gin.Context)
	FilterCategory(*gin.Context)
	DeleteCategoryImage(*gin.Context)
	AddPostInfoBlock(ctx *gin.Context)
}

func NewCategoryController(
	categoriesService *service.CategoriesService,
	categoryService *service.CategoryService,
	api *api.Api,
) CategoryController {
	return &categoryController{
		categoriesService: categoriesService,
		categoryService:   categoryService,
		api:               api,
	}
}

type categoryController struct {
	*app.BaseAjax

	categoriesService *service.CategoriesService
	categoryService   *service.CategoryService
	api               *api.Api
}

func (c *categoryController) templates(ctx *gin.Context) []contract.Template {
	templates, err := c.api.Template.GetForResources(&models.PostCategory{})
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	return templates
}
