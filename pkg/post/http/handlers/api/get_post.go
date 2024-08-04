package api

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/post/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPost(c *gin.Context) {
	id := c.Param("id")
	h := db.GetDB()
	var post models.Post

	if result := h.First(&post, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	c.JSON(http.StatusOK, &post)
}
