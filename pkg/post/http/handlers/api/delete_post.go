package api

import (
	"github.com/axlle-com/blog/pkg/common/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) DeletePost(c *gin.Context) {
	id := c.Param("id")

	var post models.Post

	if result := h.DB.First(&post, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	h.DB.Delete(&post)

	c.Status(http.StatusOK)
}
