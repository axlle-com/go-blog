package web

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	h := &handler{
		DB: db,
	}

	routes := r.Group("/admin/posts")
	routes.POST("/", h.AddPost)
	routes.GET("/", h.GetPosts)
	routes.GET("/:id", h.GetPost)
	routes.PUT("/:id", h.UpdatePost)
	routes.DELETE("/:id", h.DeletePost)

	main := r.Group("/posts")
	main.GET("/", h.GetPosts)
	main.GET("/:id", h.GetPost)
}
