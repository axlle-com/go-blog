package web

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/gallery/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	galleryRepo := repository.NewGalleryRepository()
	gallery, err := galleryRepo.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"title":   "404 Not Found",
			"content": "errors.404.gohtml",
		})
		return
	}

	imageRepo := repository.NewGalleryImageRepository()
	image, err := imageRepo.GetByID(imageId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"title":   "404 Not Found",
			"content": "errors.404.gohtml",
		})
		return
	}

	db := db.GetDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		err = imageRepo.Delete(image.ID)
		if err != nil {
			logger.New().Error(err)
			return err
		}

		err = os.Remove(image.File)
		if err != nil {
			logger.New().Error(err)
			return err
		}

		count := imageRepo.CountForGallery(gallery.ID)
		if count == 0 {
			err := galleryRepo.Delete(gallery.ID)
			if err != nil {
				logger.New().Error(err)
				return err
			}
		}
		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete file: " + err.Error(),
		})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Изображение удалено",
		})
		return
	}
}
