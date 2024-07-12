package web

import (
	. "github.com/axlle-com/blog/pkg/post"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h handler) GetPost(c *gin.Context) {
	id := c.Param("id")

	var post Post

	if result := h.DB.First(&post, id); result.Error != nil {
		c.HTML(http.StatusNotFound, "admin.404", gin.H{
			"title":   "404 Not Found",
			"content": "errors/404.gohtml",
		})
		return
	}

	c.HTML(
		http.StatusOK,
		"index",
		gin.H{
			"title":   "Home Page",
			"payload": post,
		},
	)
}
