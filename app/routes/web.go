package routes

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/axlle-com/blog/app"
	"github.com/axlle-com/blog/app/middleware"
	analyticMiddleware "github.com/axlle-com/blog/app/middleware/analytic"
	menu "github.com/axlle-com/blog/pkg/menu/models"
)

func InitWebRoutes(r *gin.Engine, container *app.Container) {
	postFrontWebController := container.PostFrontWebController()
	postController := container.PostController()
	postWebController := container.PostWebController()
	postTagWebController := container.PostTagWebController()
	postTagAjaxController := container.PostTagAjaxController()
	postCategoryWebController := container.CategoryWebController()
	postCategoryController := container.CategoryController()
	galleryController := container.GalleryAjaxController()

	fileController := container.FileController()

	userController := container.UserFrontController()

	infoBlockController := container.InfoBlockWebController()
	infoBlockAjaxController := container.InfoBlockController()

	templateController := container.TemplateWebController()
	templateAjaxController := container.TemplateController()

	messageController := container.MessageController()
	messageAjaxController := container.MessageAjaxController()
	messageFrontController := container.MessageFrontController()

	menuController := container.MenuController()
	menuItemController := container.MenuItemAjaxController()
	menuAjaxController := container.MenuAjaxController()

	analytic := analyticMiddleware.NewAnalytic(container.Queue)

	r.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})
	r.Use(middleware.Main())
	r.Use(middleware.Language())
	r.Use(middleware.Error())
	r.Use(analytic.Handler())
	r.GET("/", postFrontWebController.GetHome)
	r.GET("/login", userController.Login)
	r.POST("/auth", userController.Auth)
	r.POST("/messages", messageFrontController.CreateMessage)

	protected := r.Group("/admin")
	protected.Use(middleware.Admin())
	{
		protected.GET("", userController.Index)
		protected.GET("/logout", userController.Logout)

		protected.POST("/file/image", fileController.UploadImage)
		protected.POST("/file/images", fileController.UploadImages)
		protected.DELETE("/file/image/*filePath", fileController.DeleteImage)

		protected.GET("/posts", postWebController.GetPosts)
		protected.GET("/posts/:id", postWebController.GetPost)
		protected.GET("/posts/form", postWebController.CreatePost)
		protected.POST("/posts", postController.CreatePost)
		protected.GET("/posts/filter", postController.FilterPosts)
		protected.PUT("/posts/:id", postController.UpdatePost)
		protected.DELETE("/posts/:id", postController.DeletePost)
		protected.DELETE("/posts/:id/image", postController.DeletePostImage)
		protected.POST("/posts/:id/info-blocks/:info_block_id", postController.AddPostInfoBlock)

		protected.GET("/post/categories", postCategoryWebController.GetCategories)
		protected.GET("/post/categories/:id", postCategoryWebController.GetCategory)
		protected.POST("/post/categories", postCategoryController.CreateCategory)
		protected.PUT("/post/categories/:id", postCategoryController.UpdateCategory)
		protected.DELETE("/post/categories/:id/image", postCategoryController.DeleteCategoryImage)
		protected.POST("/post/categories/:id/info-blocks/:info_block_id", postCategoryController.AddPostInfoBlock)
		protected.DELETE("/post/categories/:id", postCategoryController.DeleteCategory)
		protected.GET("/post/categories/form", postCategoryWebController.CreateCategory)
		protected.GET("/post/categories/filter", postCategoryController.FilterCategory)

		protected.GET("/post/tags", postTagWebController.GetTags)
		protected.GET("/post/tags/form", postTagWebController.CreateTag)
		protected.GET("/post/tags/:id", postTagWebController.GetTag)
		protected.POST("/post/tags", postTagAjaxController.Create)
		protected.PUT("/post/tags/:id", postTagAjaxController.Update)
		protected.DELETE("/post/tags/:id", postTagAjaxController.Delete)
		protected.DELETE("/post/tags/:id/image", postTagAjaxController.DeleteImage)
		protected.GET("/post/tags/filter", postTagAjaxController.Filter)

		protected.GET("/info-blocks", infoBlockController.GetInfoBlocks)
		protected.GET("/info-blocks/:id", infoBlockController.GetInfoBlock)
		protected.GET("/info-blocks/form", infoBlockController.CreateInfoBlock)

		protected.POST("/info-blocks", infoBlockAjaxController.CreateInfoBlock)
		protected.PUT("/info-blocks/:id", infoBlockAjaxController.UpdateInfoBlock)
		protected.DELETE("/info-blocks/:id", infoBlockAjaxController.DeleteInfoBlock)
		protected.DELETE("/info-blocks/:id/image", infoBlockAjaxController.DeleteBlockImage)
		protected.GET("/info-blocks/filter", infoBlockAjaxController.FilterInfoBlock)
		protected.GET("/ajax/info-blocks/:id", infoBlockAjaxController.GetInfoBlock)
		protected.GET("/ajax/info-blocks/:id/card", infoBlockAjaxController.GetInfoBlockCard)
		protected.DELETE("/ajax/info-blocks/:id/detach", infoBlockAjaxController.DetachInfoBlock)

		protected.GET("/templates", templateController.GetTemplates)
		protected.GET("/templates/:id", templateController.GetTemplate)

		protected.POST("/templates", templateAjaxController.CreateTemplate)
		protected.PUT("/templates/:id", templateAjaxController.UpdateTemplate)
		protected.DELETE("/templates/:id", templateAjaxController.DeleteTemplate)
		protected.GET("/templates/filter", templateAjaxController.FilterTemplate)
		protected.GET("/templates/resources/:template", templateAjaxController.GetResourceTemplate)

		protected.GET("/messages", messageController.GetMessages)
		protected.GET("/messages/:id", messageController.GetMessage)
		protected.GET("/ajax/messages", messageAjaxController.GetMessages)
		protected.GET("/ajax/messages/:id", messageAjaxController.GetMessage)
		protected.DELETE("/ajax/messages/:id", messageAjaxController.DeleteMessage)

		protected.DELETE("/gallery/:id/image/:image_id", galleryController.DeleteImage)

		protected.GET("/menus", menuController.GetMenus)
		protected.GET("/menus/form", menuController.CreateMenu)
		protected.GET("/menus/:id", menuController.GetMenu)
		protected.GET("/ajax/menus/menus-items", menuItemController.GetMenuItems)
		protected.POST("/menus", menuAjaxController.Create)
		protected.PUT("/menus/:id", menuAjaxController.Update)
	}

	r.GET("/:alias", postFrontWebController.GetPost)

	r.NoRoute(func(ctx *gin.Context) {
		path := ctx.Request.URL.Path

		if strings.HasPrefix(path, "/admin") {
			ctx.HTML(http.StatusNotFound, "admin.404", gin.H{
				"title": "Админка — 404",
				"menu":  menu.NewMenu(ctx.FullPath()),
			})
		} else {
			ctx.HTML(http.StatusNotFound, "404", gin.H{
				"title": "Страница не найдена",
			})
		}
	})
}
