package api

import (
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) GetPost(ctx *gin.Context) {
	var post models.Post
	ctx.JSON(http.StatusOK, &post)
}
