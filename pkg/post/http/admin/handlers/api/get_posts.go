package api

import (
	"github.com/axlle-com/blog/pkg/post/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) GetPosts(ctx *gin.Context) {
	var posts []models.Post
	ctx.JSON(http.StatusOK, &posts)
}
