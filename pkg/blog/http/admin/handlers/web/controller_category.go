package web

import (
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/blog/service"
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	"github.com/axlle-com/blog/pkg/info_block/provider"
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
	infoBlockProvider provider.InfoBlockProvider,
) ControllerCategory {
	return &controllerCategory{
		categoriesService: categoriesService,
		categoryService:   categoryService,
		templateProvider:  template,
		userProvider:      user,
		galleryProvider:   gallery,
		infoBlockProvider: infoBlockProvider,
	}
}

type controllerCategory struct {
	*app.BaseAjax

	categoriesService *service.CategoriesService
	categoryService   *service.CategoryService
	templateProvider  template.TemplateProvider
	userProvider      user.UserProvider
	galleryProvider   gallery.GalleryProvider
	infoBlockProvider provider.InfoBlockProvider
}
