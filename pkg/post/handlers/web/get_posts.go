package web

import (
	. "github.com/axlle-com/blog/pkg/post"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h handler) GetPosts(c *gin.Context) {
	var posts []Post

	if result := h.DB.Find(&posts); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	c.HTML(
		http.StatusOK,
		"index",
		gin.H{
			"title":   "Home Page",
			"payload": posts,
		},
	)
}
