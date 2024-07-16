package api

import (
	"github.com/axlle-com/blog/pkg/common/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) GetPosts(c *gin.Context) {
	var posts []models.Post

	if result := h.DB.Find(&posts); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	c.JSON(http.StatusOK, &posts)
}
