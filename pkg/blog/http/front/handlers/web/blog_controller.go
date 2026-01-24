package web

import (
	"html/template"
	"strings"

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
	config contract.Config,
	api *api.Api,
	view contract.View,
	service *service.PostService,
	services *service.PostCollectionService,
	category *service.CategoryService,
	categories *service.CategoriesService,
) BlogController {
	return &blogController{
		config:                config,
		api:                   api,
		view:                  view,
		postService:           service,
		postCollectionService: services,
		categoryService:       category,
		categoriesService:     categories,
	}
}

type blogController struct {
	*app.BaseAjax

	config                contract.Config
	api                   *api.Api
	view                  contract.View
	postService           *service.PostService
	postCollectionService *service.PostCollectionService
	categoryService       *service.CategoryService
	categoriesService     *service.CategoriesService
}

func (c *blogController) settings(ctx *gin.Context, publisher contract.Publisher) map[string]any {
	menu, err := c.api.Menu.GetMenuString(1, ctx.Request.URL.Path)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
	}

	appHost := strings.TrimRight(c.config.AppHost(), "/")

	var pubURL, pubImage, pubTitle, pubDesc string
	if publisher != nil {
		pubURL = publisher.GetURL()
		pubTitle = publisher.GetMetaTitle()
		pubDesc = publisher.GetMetaDescription()

		img := publisher.GetImage()
		if img != "" && strings.HasPrefix(img, "/") {
			img = joinURL(c.config.AppHost(), img)
		}
		pubImage = img
	}

	base := joinURL(appHost, pubURL)

	body := gin.H{
		"menu":      template.HTML(menu),
		"user":      c.GetAdmin(ctx),
		"csrfToken": csrf.GetToken(ctx),

		"baseURL":         appHost,
		"metaURL":         base,
		"metaImage":       pubImage,
		"metaTitle":       pubTitle,
		"metaDescription": pubDesc,
	}

	company, ok := c.api.CompanyInfo.GetCompanyInfo(c.config.Layout(), "global")
	if company != nil && ok {
		body["company"] = company
	}

	return body
}

func joinURL(base, path string) string {
	base = strings.TrimRight(base, "/")
	path = strings.TrimLeft(path, "/")

	if path == "" {
		return base
	}
	if base == "" {
		return "/" + path
	}

	return base + "/" + path
}
