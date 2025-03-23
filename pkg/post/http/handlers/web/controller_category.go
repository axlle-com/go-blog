package web

import (
	app "github.com/axlle-com/blog/pkg/app/models"
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	"github.com/axlle-com/blog/pkg/post/repository"
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
	categoryRepo repository.CategoryRepository,
	categoryService *service.CategoriesService,
	template template.TemplateProvider,
	user user.UserProvider,
	gallery gallery.GalleryProvider,
) ControllerCategory {
	return &controllerCategory{
		categoryRepo:    categoryRepo,
		categoryService: categoryService,
		template:        template,
		user:            user,
		gallery:         gallery,
	}
}

type controllerCategory struct {
	*app.BaseAjax

	categoryRepo    repository.CategoryRepository
	categoryService *service.CategoriesService
	template        template.TemplateProvider
	user            user.UserProvider
	gallery         gallery.GalleryProvider
}
