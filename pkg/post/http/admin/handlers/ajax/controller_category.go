package ajax

import (
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/post/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
)

type CategoryController interface {
	UpdateCategory(*gin.Context)
	CreateCategory(*gin.Context)
	DeleteCategory(*gin.Context)
	FilterCategory(*gin.Context)
	DeleteCategoryImage(*gin.Context)
}

func NewCategoryController(
	categoriesService *service.CategoriesService,
	categoryService *service.CategoryService,
	template template.TemplateProvider,
	user user.UserProvider,
) CategoryController {
	return &categoryController{
		categoriesService: categoriesService,
		categoryService:   categoryService,
		templateProvider:  template,
		userProvider:      user,
	}
}

type categoryController struct {
	*app.BaseAjax

	categoriesService *service.CategoriesService
	categoryService   *service.CategoryService
	templateProvider  template.TemplateProvider
	userProvider      user.UserProvider
}
