package routes

import (
	post "github.com/axlle-com/blog/pkg/post/http/handlers/web"
	user "github.com/axlle-com/blog/pkg/user/http/handlers/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitializeWebRoutes(r *gin.Engine) {
	r.GET("/", ShowIndexPage)
	post.RegisterRoutes(r)
	user.RegisterRoutes(r)
}

func ShowIndexPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"index",
		gin.H{
			"title":   "Home Page",
			"payload": nil,
		},
	)
}
