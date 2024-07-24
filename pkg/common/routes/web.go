package routes

import (
	"github.com/axlle-com/blog/pkg/common/middleware"
	post "github.com/axlle-com/blog/pkg/post/http/handlers/web"
	user "github.com/axlle-com/blog/pkg/user/http/handlers/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitializeWebRoutes(r *gin.Engine) {
	r.GET("/", ShowIndexPage)
	r.GET("/login", user.Login)
	r.POST("/auth", user.Auth)
	r.POST("/user", user.CreateUser)
	r.GET("/posts", post.GetPosts)

	protected := r.Group("/admin")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/", user.Index)
		protected.POST("/posts", post.CreatePost)
		protected.GET("/posts", post.GetPosts)
		protected.GET("/posts/:id", post.GetPost)
		protected.PUT("/posts/:id", post.UpdatePost)
		protected.DELETE("/posts/:id", post.DeletePost)
	}
	r.GET("/:alias", post.GetPost)
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
