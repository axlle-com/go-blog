package api

import (
	"github.com/axlle-com/blog/pkg/app/db"
	"github.com/axlle-com/blog/pkg/post/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) GetPosts(ctx *gin.Context) {
	var posts []models.Post
	h := db.GetDB()
	if result := h.Find(&posts); result.Error != nil {
		ctx.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, &posts)
}
