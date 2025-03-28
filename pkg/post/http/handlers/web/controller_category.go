package web

import (
	app "github.com/axlle-com/blog/pkg/app/models"
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	"github.com/axlle-com/blog/pkg/post/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/gin-gonic/gin"
)

type ControllerCategory interface {
	GetCategory(*gin.Context)
	GetCategories(*gin.Context)
	CreateCategory(*gin.Context)
}

func NewWebControllerCategory(
	categoriesService *service.CategoriesService,
	categoryService *service.CategoryService,
	template template.TemplateProvider,
	user user.UserProvider,
	gallery gallery.GalleryProvider,
) ControllerCategory {
	return &controllerCategory{
		categoriesService: categoriesService,
		categoryService:   categoryService,
		templateProvider:  template,
		userProvider:      user,
		galleryProvider:   gallery,
	}
}

type controllerCategory struct {
	*app.BaseAjax

	categoriesService *service.CategoriesService
	categoryService   *service.CategoryService
	templateProvider  template.TemplateProvider
	userProvider      user.UserProvider
	galleryProvider   gallery.GalleryProvider
}
