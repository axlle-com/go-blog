package web

import (
	"github.com/axlle-com/blog/pkg/app/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) DeleteImage(ctx *gin.Context) {
	id := c.GetID(ctx)
	if id == 0 {
		return
	}

	imageId := c.getImageID(ctx)
	if imageId == 0 {
		return
	}

	_, err := c.gallery.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		ctx.Abort()
		return
	}

	image, err := c.image.GetByID(imageId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		ctx.Abort()
		return
	}

	err = c.imageService.DeleteImage(image)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete file: " + err.Error(),
		})
		ctx.Abort()
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "Изображение удалено"})
		return
	}
}
