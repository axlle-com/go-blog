package routes

import (
	"github.com/axlle-com/blog/pkg/app"
	"github.com/gin-gonic/gin"
)

func InitializeApiRoutes(r *gin.Engine, container *app.Container) {
	controller := container.PostApiController()

	r.POST("/api/posts", controller.CreatePost)
	r.GET("/api/posts", controller.GetPosts)
	r.GET("/api/posts/:id", controller.GetPost)
	r.PUT("/api/posts/:id", controller.UpdatePost)
	r.DELETE("/api/posts/:id", controller.DeletePost)
}
