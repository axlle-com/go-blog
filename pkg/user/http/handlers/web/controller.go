package web

import (
	"github.com/axlle-com/blog/pkg/common/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/login", Login)
	r.POST("/auth", Auth)
	r.POST("/user", CreateUser)

	protected := r.Group("/admin")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/", Index)
	}
}
