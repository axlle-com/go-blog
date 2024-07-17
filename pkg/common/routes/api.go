package routes

import (
	post "github.com/axlle-com/blog/pkg/post/http/handlers/api"
	"github.com/gin-gonic/gin"
)

func InitializeApiRoutes(r *gin.Engine) {
	post.RegisterRoutes(r)
}
