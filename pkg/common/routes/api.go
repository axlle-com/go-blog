package routes

import (
	post "github.com/axlle-com/blog/pkg/post/http/handlers/api"
	"github.com/gin-gonic/gin"
)

func InitializeApiRoutes(r *gin.Engine) {
	r.POST("/api/posts", post.AddPost)
	r.GET("/api/posts", post.GetPosts)
	r.GET("/api/posts/:id", post.GetPost)
	r.PUT("/api/posts/:id", post.UpdatePost)
	r.DELETE("/api/posts/:id", post.DeletePost)
}
