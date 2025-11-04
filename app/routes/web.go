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

	r.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})
	r.Use(middleware.Main())
	r.Use(middleware.Language(container.I18n))
	r.Use(middleware.Error())
	r.Use(analytic.Handler())
	r.GET("/", container.FrontWebPostController.GetHome)
	r.GET("/login", container.FrontWebUserController.Login)
	r.POST("/auth", container.FrontWebUserController.Auth)
	r.POST("/messages", container.FrontAjaxMessageController.CreateMessage)

	protected := r.Group("/admin")
	protected.Use(middleware.Admin())
	{
		protected.GET("", container.FrontWebUserController.Index)
		protected.GET("/logout", container.FrontWebUserController.Logout)

		protected.POST("/file/image", container.AdminWebFileController.UploadImage)
		protected.POST("/file/images", container.AdminWebFileController.UploadImages)
		protected.DELETE("/file/image/*filePath", container.AdminWebFileController.DeleteImage)

		protected.GET("/posts", container.AdminWebPostController.GetPosts)
		protected.GET("/posts/:id", container.AdminWebPostController.GetPost)
		protected.GET("/posts/form", container.AdminWebPostController.CreatePost)
		protected.POST("/posts", container.AdminAjaxPostController.CreatePost)
		protected.GET("/posts/filter", container.AdminAjaxPostController.FilterPosts)
		protected.PUT("/posts/:id", container.AdminAjaxPostController.UpdatePost)
		protected.DELETE("/posts/:id", container.AdminAjaxPostController.DeletePost)
		protected.DELETE("/posts/:id/image", container.AdminAjaxPostController.DeletePostImage)
		protected.POST("/posts/:id/info-blocks/:info_block_id", container.AdminAjaxPostController.AddPostInfoBlock)

		protected.GET("/post/categories", container.AdminWebCategoryController.GetCategories)
		protected.GET("/post/categories/:id", container.AdminWebCategoryController.GetCategory)
		protected.POST("/post/categories", container.AdminAjaxCategoryController.CreateCategory)
		protected.PUT("/post/categories/:id", container.AdminAjaxCategoryController.UpdateCategory)
		protected.DELETE("/post/categories/:id/image", container.AdminAjaxCategoryController.DeleteCategoryImage)
		protected.POST("/post/categories/:id/info-blocks/:info_block_id", container.AdminAjaxCategoryController.AddPostInfoBlock)
		protected.DELETE("/post/categories/:id", container.AdminAjaxCategoryController.DeleteCategory)
		protected.GET("/post/categories/form", container.AdminWebCategoryController.CreateCategory)
		protected.GET("/post/categories/filter", container.AdminAjaxCategoryController.FilterCategory)

		protected.GET("/post/tags", container.AdminWebTagController.GetTags)
		protected.GET("/post/tags/form", container.AdminWebTagController.CreateTag)
		protected.GET("/post/tags/:id", container.AdminWebTagController.GetTag)
		protected.POST("/post/tags", container.AdminAjaxTagController.Create)
		protected.PUT("/post/tags/:id", container.AdminAjaxTagController.Update)
		protected.DELETE("/post/tags/:id", container.AdminAjaxTagController.Delete)
		protected.DELETE("/post/tags/:id/image", container.AdminAjaxTagController.DeleteImage)
		protected.GET("/post/tags/filter", container.AdminAjaxTagController.Filter)

		protected.GET("/info-blocks", container.AdminWebInfoBlockController.GetInfoBlocks)
		protected.GET("/info-blocks/:id", container.AdminWebInfoBlockController.GetInfoBlock)
		protected.GET("/info-blocks/form", container.AdminWebInfoBlockController.CreateInfoBlock)

		protected.POST("/info-blocks", container.AdminAjaxInfoBlockController.CreateInfoBlock)
		protected.PUT("/info-blocks/:id", container.AdminAjaxInfoBlockController.UpdateInfoBlock)
		protected.DELETE("/info-blocks/:id", container.AdminAjaxInfoBlockController.DeleteInfoBlock)
		protected.DELETE("/info-blocks/:id/image", container.AdminAjaxInfoBlockController.DeleteBlockImage)
		protected.GET("/info-blocks/filter", container.AdminAjaxInfoBlockController.FilterInfoBlock)
		protected.GET("/ajax/info-blocks/:id", container.AdminAjaxInfoBlockController.GetInfoBlock)
		protected.GET("/ajax/info-blocks/:id/card", container.AdminAjaxInfoBlockController.GetInfoBlockCard)
		protected.DELETE("/ajax/info-blocks/:id/detach", container.AdminAjaxInfoBlockController.DetachInfoBlock)

		protected.GET("/templates", container.AdminWebTemplateController.GetTemplates)
		protected.GET("/templates/:id", container.AdminWebTemplateController.GetTemplate)

		protected.POST("/templates", container.AdminAjaxTemplateController.CreateTemplate)
		protected.PUT("/templates/:id", container.AdminAjaxTemplateController.UpdateTemplate)
		protected.DELETE("/templates/:id", container.AdminAjaxTemplateController.DeleteTemplate)
		protected.GET("/templates/filter", container.AdminAjaxTemplateController.FilterTemplate)
		protected.GET("/templates/resources/:template", container.AdminAjaxTemplateController.GetResourceTemplate)

		protected.GET("/messages", container.AdminWebMessageController.GetMessages)
		protected.GET("/messages/:id", container.AdminWebMessageController.GetMessage)
		protected.GET("/ajax/messages", container.AdminAjaxMessageController.GetMessages)
		protected.GET("/ajax/messages/:id", container.AdminAjaxMessageController.GetMessage)
		protected.DELETE("/ajax/messages/:id", container.AdminAjaxMessageController.DeleteMessage)

		protected.DELETE("/gallery/:id/image/:image_id", container.AdminAjaxGalleryController.DeleteImage)

		protected.GET("/menus", container.AdminWebMenuController.GetMenus)
		protected.GET("/menus/form", container.AdminWebMenuController.CreateMenu)
		protected.GET("/menus/:id", container.AdminWebMenuController.GetMenu)
		protected.GET("/ajax/menus/menus-items", container.AdminAjaxMenuItemController.GetMenuItems)
		protected.POST("/menus", container.AdminAjaxMenuController.Create)
		protected.PUT("/menus/:id", container.AdminAjaxMenuController.Update)
	}

	r.GET("/:alias", container.FrontWebPostController.GetPost)

	r.NoRoute(func(ctx *gin.Context) {
		path := ctx.Request.URL.Path

		if strings.HasPrefix(path, "/admin") {
			ctx.HTML(http.StatusNotFound, "admin.404", gin.H{
				"title": "Админка — 404",
				"menu":  menu.NewMenu(ctx.FullPath(), nil),
			})
		} else {
			ctx.HTML(http.StatusNotFound, "404", gin.H{
				"title": "Страница не найдена",
			})
		}
	})
}
