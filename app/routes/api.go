package routes

import (
	"github.com/axlle-com/blog/app"
	"github.com/gin-gonic/gin"
)

func InitApiRoutes(r *gin.Engine, container *app.Container) {
	controller := container.AdminApiPostController

	r.POST("/api/posts", controller.CreatePost)
	r.GET("/api/posts", controller.GetPosts)
	r.GET("/api/posts/:id", controller.GetPost)
	r.PUT("/api/posts/:id", controller.UpdatePost)
	r.DELETE("/api/posts/:id", controller.DeletePost)
}
