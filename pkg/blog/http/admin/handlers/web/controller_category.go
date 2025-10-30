package web

import (
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	appPovider "github.com/axlle-com/blog/app/models/provider"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/axlle-com/blog/pkg/blog/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
)

type CategoryController interface {
	GetCategory(*gin.Context)
	GetCategories(*gin.Context)
	CreateCategory(*gin.Context)
}

func NewWebCategoryController(
	categoriesService *service.CategoriesService,
	categoryService *service.CategoryService,
	template template.TemplateProvider,
	user user.UserProvider,
	gallery appPovider.GalleryProvider,
	infoBlockProvider appPovider.InfoBlockProvider,
) CategoryController {
	return &categoryController{
		categoriesService: categoriesService,
		categoryService:   categoryService,
		templateProvider:  template,
		userProvider:      user,
		galleryProvider:   gallery,
		infoBlockProvider: infoBlockProvider,
	}
}

type categoryController struct {
	*app.BaseAjax

	categoriesService *service.CategoriesService
	categoryService   *service.CategoryService
	templateProvider  template.TemplateProvider
	userProvider      user.UserProvider
	galleryProvider   appPovider.GalleryProvider
	infoBlockProvider appPovider.InfoBlockProvider
}

func (c *categoryController) templates(ctx *gin.Context) []contract.Template {
	templates, err := c.templateProvider.GetForResources(&models.PostCategory{})
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	return templates
}
