package routes

import (
	"github.com/axlle-com/blog/pkg/common/middleware"
	gallery "github.com/axlle-com/blog/pkg/gallery/http/handlers/web"
	postAjax "github.com/axlle-com/blog/pkg/post/http/handlers/ajax"
	post "github.com/axlle-com/blog/pkg/post/http/handlers/web"
	user "github.com/axlle-com/blog/pkg/user/http/handlers/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitializeWebRoutes(r *gin.Engine) {
	postController := postAjax.NewController(r)
	postWebController := post.NewWebController(r)
	galleryController := gallery.NewController(r)

	r.GET("/", ShowIndexPage)
	r.GET("/login", user.Login)
	r.POST("/auth", user.Auth)
	r.POST("/user", user.CreateUser)

	protected := r.Group("/admin")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/", user.Index)
		protected.GET("/logout", user.Logout)
		protected.GET("/posts", postWebController.GetPosts)
		protected.GET("/posts/:id", postWebController.GetPost)
		protected.GET("/posts/form", postWebController.CreatePost)

		protected.POST("/posts", postController.CreatePost)
		protected.POST("/posts/filter", postController.FilterPosts)
		protected.PUT("/posts/:id", postController.UpdatePost)
		protected.DELETE("/posts/:id", postController.DeletePost)
		protected.DELETE("/posts/:id/image", postController.DeletePostImage)
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
