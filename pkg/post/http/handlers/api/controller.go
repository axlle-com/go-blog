package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	routes := r.Group("/api/posts")
	routes.POST("/", AddPost)
	routes.GET("/", GetPosts)
	routes.GET("/:id", GetPost)
	routes.PUT("/:id", UpdatePost)
	routes.DELETE("/:id", DeletePost)
}
