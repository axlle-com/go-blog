package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Controller interface {
	DeleteImage(*gin.Context)
}

func NewController(r *gin.Engine) Controller {
	return &controller{engine: r}
}

type controller struct {
	engine *gin.Engine
}

func (c *controller) GetID(ctx *gin.Context) uint {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		ctx.Abort()
	}
	return uint(id)
}

func (c *controller) getImageID(ctx *gin.Context) uint {
	idParam := ctx.Param("image_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		ctx.Abort()
	}
	return uint(id)
}
