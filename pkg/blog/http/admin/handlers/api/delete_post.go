package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) DeletePost(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
