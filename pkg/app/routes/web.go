package routes

import (
	"fmt"
	models2 "github.com/axlle-com/blog/pkg/menu/models"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"net/http"
	"strings"

	"github.com/axlle-com/blog/pkg/app"
	"github.com/axlle-com/blog/pkg/app/middleware"
	file "github.com/axlle-com/blog/pkg/file/http"
	post "github.com/axlle-com/blog/pkg/post/http/handlers/web"
	user "github.com/axlle-com/blog/pkg/user/http/handlers/web"
)

func InitializeWebRoutes(r *gin.Engine, container *app.Container) {
	postController := container.PostController()
	postWebController := container.PostWebController()
	postCategoryWebController := container.CategoryWebController()
	postCategoryController := container.CategoryController()
	galleryController := container.GalleryAjaxController()

	fileController := file.NewFileController(
		container.FileService,
	)

	userController := user.NewUserWebController(
		container.UserService,
		container.UserAuthService,
	)

	infoBlockController := container.InfoBlockWebController()
	infoBlockAjaxController := container.InfoBlockController()

	r.Use(middleware.Error())
	r.GET("/", ShowIndexPage)
	r.GET("/login", userController.Login)
	r.POST("/auth", userController.Auth)
	r.POST("/user", userController.CreateUser)

	protected := r.Group("/admin")
	protected.Use(middleware.AuthRequired())
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

		protected.GET("/categories", postCategoryWebController.GetCategories)
		protected.GET("/categories/:id", postCategoryWebController.GetCategory)
		protected.POST("/categories", postCategoryController.CreateCategory)
		protected.PUT("/categories/:id", postCategoryController.UpdateCategory)
		protected.DELETE("/categories/:id/image", postCategoryController.DeleteCategoryImage)
		protected.DELETE("/categories/:id", postCategoryController.DeleteCategory)
		protected.GET("/categories/form", postCategoryWebController.CreateCategory)
		protected.GET("/categories/filter", postCategoryController.FilterCategory)

		protected.GET("/info-blocks", infoBlockController.GetInfoBlocks)
		protected.GET("/info-blocks/:id", infoBlockController.GetInfoBlock)

		protected.POST("/info-blocks", infoBlockAjaxController.CreateInfoBlock)
		protected.PUT("/info-blocks/:id", infoBlockAjaxController.UpdateInfoBlock)
		protected.DELETE("/info-blocks/:id", infoBlockAjaxController.DeleteInfoBlock)
		protected.GET("/info-blocks/filter", infoBlockAjaxController.FilterInfoBlock)

		protected.DELETE("/gallery/:id/image/:image_id", galleryController.DeleteImage)
	}
	r.GET("/:alias", post.GetPostFront)

	r.NoRoute(func(ctx *gin.Context) {
		path := ctx.Request.URL.Path

		if strings.HasPrefix(path, "/admin") {
			ctx.HTML(http.StatusNotFound, "admin.404", gin.H{
				"title": "Админка — 404",
				"menu":  models2.NewMenu(ctx.FullPath()),
			})
		} else {
			ctx.HTML(http.StatusNotFound, "404", gin.H{
				"title": "Страница не найдена",
			})
		}
	})
}

func ShowIndexPage(c *gin.Context) {
	ua := user_agent.New(c.GetHeader("User-Agent"))
	name, version := ua.Browser()

	fmt.Println("Browser:", name, version)
	fmt.Println("OS:", ua.OS())
	fmt.Println("Mobile?", ua.Mobile())
	fmt.Println(DetectDeviceType(c.GetHeader("User-Agent")))

	c.HTML(
		http.StatusOK,
		"index",
		gin.H{
			"title":   "Home Page",
			"payload": nil,
		},
	)
}

func DetectDeviceType(uaString string) string {
	ua := user_agent.New(uaString)

	switch {
	case ua.Bot():
		return "bot"
	case strings.Contains(uaString, "iPad"):
		return "tablet"
	case ua.Mobile():
		return "mobile"
	default:
		return "desktop"
	}
}
