package web

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func (c *controller) DeleteImage(ctx *gin.Context) {
	id := c.getID(ctx)
	if id == 0 {
		return
	}

	imageId := c.getImageID(ctx)
	if imageId == 0 {
		return
	}

	galleryRepo := models.NewGalleryRepository()
	gallery, err := galleryRepo.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	imageRepo := models.NewGalleryImageRepository()
	image, err := imageRepo.GetByID(imageId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Ресурс не найден"})
		return
	}

	err = imageRepo.Delete(image.ID)
	if err != nil {
		logger.New().Error(err)
	}

	err = os.Remove(image.GetFilePath())
	if err != nil {
		logger.New().Error(err)
	}

	count := imageRepo.CountForGallery(gallery.ID)
	if count == 0 {
		err := galleryRepo.Delete(gallery.ID)
		if err != nil {
			logger.New().Error(err)
		}
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete file: " + err.Error(),
		})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "Изображение удалено"})
		return
	}
}
