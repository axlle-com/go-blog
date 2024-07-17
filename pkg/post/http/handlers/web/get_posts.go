package web

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPosts(c *gin.Context) {
	var posts []models.Post
	h := db.GetDB()
	if result := h.Find(&posts); result.Error != nil {
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
