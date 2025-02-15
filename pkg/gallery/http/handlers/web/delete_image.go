package web

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"github.com/axlle-com/blog/pkg/gallery/service"
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

	galleryRepo := models.GalleryRepo()
	_, err := galleryRepo.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		ctx.Abort()
		return
	}

	imageRepo := models.ImageRepo()
	image, err := imageRepo.GetByID(imageId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		ctx.Abort()
		return
	}

	err = service.DeleteImage(image)
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
