package routes

import (
	"github.com/axlle-com/blog/pkg/app"
	"github.com/axlle-com/blog/pkg/app/middleware"
	file "github.com/axlle-com/blog/pkg/file/http"
	post "github.com/axlle-com/blog/pkg/post/http/handlers/web"
	user "github.com/axlle-com/blog/pkg/user/http/handlers/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitializeWebRoutes(r *gin.Engine, container *app.Container) {
	postController := container.PostController()
	postWebController := container.PostWebController()
	postCategoryWebController := container.PostCategoryWebController()
	galleryController := container.GalleryAjaxController()

	fileController := file.New(
		container.FileService,
	)

	userController := user.New(
		container.UserRepo,
	)

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
		protected.GET("/categories/form", postCategoryWebController.CreateCategory)

		protected.DELETE("/gallery/:id/image/:image_id", galleryController.DeleteImage)
	}
	r.GET("/:alias", post.GetPostFront)
}

func ShowIndexPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"index",
		gin.H{
			"title":   "Home Page",
			"payload": nil,
		},
	)
}
