package web

import (
	"github.com/axlle-com/blog/pkg/menu"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) Index(ctx *gin.Context) {
	ctx.HTML(
		http.StatusOK,
		"admin.index",
		gin.H{
			"title": "dashboard",
			"menu":  menu.NewMenu(ctx.FullPath()),
		},
	)
}
