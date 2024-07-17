package web

import (
	"github.com/axlle-com/blog/pkg/common/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/:alias", GetPost)
	r.GET("/posts", GetPosts)

	protected := r.Group("/admin")
	protected.Use(middleware.AuthRequired())
	{
		protected.POST("/posts", CreatePost)
		protected.GET("/posts", GetPosts)
		protected.GET("/posts/:id", GetPost)
		protected.PUT("/posts/:id", UpdatePost)
		protected.DELETE("/posts/:id", DeletePost)
	}
}
