package api

import (
	"github.com/axlle-com/blog/pkg/post"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) GetPost(c *gin.Context) {
	id := c.Param("id")

	var post post.Post

	if result := h.DB.First(&post, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	c.JSON(http.StatusOK, &post)
}
