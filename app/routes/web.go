package routes

import (
	"net/http"
	"strings"

	"github.com/axlle-com/blog/app/di"
	"github.com/axlle-com/blog/app/middleware"
	analyticMiddleware "github.com/axlle-com/blog/app/middleware/analytic"
	"github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	menu "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func InitWebRoutes(r *gin.Engine, config contract.Config, container *di.Container) {
	analytic := analyticMiddleware.NewAnalytic(container.Queue)
	errorMiddleware := middleware.NewError(container.View)
	languageMiddleware := middleware.NewLanguage(container.I18n)

	r.Use(func(ctx *gin.Context) {
		if strings.HasPrefix(ctx.Request.URL.Path, "/.well-known/") {
			ctx.Status(http.StatusNoContent)
			ctx.Abort()
			return
		}
		ctx.Next()
	})
	r.Use(sessions.Sessions(config.SessionsName(), container.Store))
	r.Use(gzip.Gzip(gzip.BestSpeed))
	r.Use(csrf.Middleware(csrf.Options{
		Secret: string(config.KeyCookie()),
		ErrorFunc: func(ctx *gin.Context) {
			ctx.String(http.StatusForbidden, "CSRF token mismatch")
			ctx.Abort()
		},
	}))

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
	r.POST("/messages/:form", container.FrontAjaxMessageController.CreateMessage)

	InitAdminRoutes(r, container)

	r.GET("/:alias", container.FrontWebPostController.RenderByURL)

	r.NoRoute(func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		base := models.BaseAjax{}

		var tpl string
		var body gin.H

		if strings.HasPrefix(path, "/admin") {
			tpl = "admin.404"
			body = gin.H{
				"title": "Page not found",
				"error": "404",
				"menu":  menu.NewMenu(ctx.FullPath(), nil),
			}
		} else {
			tpl = container.View.ViewStatic("error")
			body = gin.H{
				"title": "Page not found",
				"error": "404",
			}
		}

		base.RenderHTML(ctx, http.StatusNotFound, tpl, body)
	})
}
