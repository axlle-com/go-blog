package web

import (
	"github.com/axlle-com/blog/pkg/menu"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"admin.index",
		gin.H{
			"title": "dashboard",
			"menu":  menu.NewMenu(c.FullPath()),
		},
	)
}
