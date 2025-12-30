package routes

import (
	"net/http"
	"strings"

	"github.com/axlle-com/blog/app/di"
	"github.com/axlle-com/blog/app/middleware"
	analyticMiddleware "github.com/axlle-com/blog/app/middleware/analytic"
	menu "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
)

func InitWebRoutes(r *gin.Engine, container *di.Container) {
	analytic := analyticMiddleware.NewAnalytic(container.Queue)
	errorMiddleware := middleware.NewError(container.View)
	languageMiddleware := middleware.NewLanguage(container.I18n)

	r.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})
	r.Use(middleware.Main())
	r.Use(languageMiddleware.Handler())
	r.Use(errorMiddleware.Handler())
	r.Use(analytic.Handler())
	r.GET("/", container.FrontWebPostController.RenderHome)
	r.GET("/login", container.FrontWebUserController.Login)
	r.POST("/auth", container.FrontWebUserController.Auth)
	r.POST("/messages", container.FrontAjaxMessageController.CreateMessage)

	InitAdminRoutes(r, container)

	r.GET("/:alias", container.FrontWebPostController.RenderByURL)

	r.NoRoute(func(ctx *gin.Context) {
		path := ctx.Request.URL.Path

		if strings.HasPrefix(path, "/admin") {
			ctx.HTML(http.StatusNotFound, "admin.404", gin.H{
				"title": "Админка — 404",
				"menu":  menu.NewMenu(ctx.FullPath(), nil),
			})
		} else {
			ctx.HTML(http.StatusNotFound, container.View.View("error"), gin.H{"title": "Page not found", "error": "404"})
		}
	})
}
