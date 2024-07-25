package web

import (
	"github.com/axlle-com/blog/pkg/common/db"
	. "github.com/axlle-com/blog/pkg/common/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetPost(c *gin.Context) {
	currentRoute := c.MustGet("currentRoute").(string)
	log.Println(currentRoute)
	id := c.Param("id")
	h := db.GetDB()
	var post Post

	if result := h.First(&post, id); result.Error != nil {
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
