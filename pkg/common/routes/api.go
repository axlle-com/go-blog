package routes

import (
	post "github.com/axlle-com/blog/pkg/post/http/handlers/api"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitializeApiRoutes(r *gin.Engine, db *gorm.DB) {
	post.RegisterRoutes(r, db)
}
