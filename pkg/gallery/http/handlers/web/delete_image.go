package web

import (
	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/logger"
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
		ctx.JSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		ctx.Abort()
		return
	}

	image, err := c.image.GetByID(imageId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": errutil.ResourceNotfound})
		ctx.Abort()
		return
	}

	err = c.imageService.DeleteImage(image)
	if err != nil {
		logger.WithRequest(ctx).Error(err)
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
