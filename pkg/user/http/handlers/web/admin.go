package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"admin.index",
		gin.H{
			"Title": "Авторизация",
		},
	)
}
