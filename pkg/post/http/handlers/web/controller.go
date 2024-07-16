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

	r.POST("/admin/posts", h.CreatePost)
	r.GET("/admin/posts", h.GetPosts)
	r.GET("/admin/posts/:id", h.GetPost)
	r.PUT("/admin/posts/:id", h.UpdatePost)
	r.DELETE("/admin/posts/:id", h.DeletePost)

	r.GET("/:alias", h.GetPost)
	r.GET("/posts", h.GetPosts)
}
