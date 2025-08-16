package api

import (
	"net/http"

	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/gin-gonic/gin"
)

func (c *controller) GetPosts(ctx *gin.Context) {
	var posts []models.Post
	ctx.JSON(http.StatusOK, &posts)
}
