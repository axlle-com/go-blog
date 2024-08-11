package service

import (
	"fmt"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"path/filepath"
)

func SaveImage(image *models.GalleryImage, c *gin.Context) error {
	if image.FileHeader != nil {
		newFileName := fmt.Sprintf("/public/uploads/%d/%s%s", image.GalleryID, uuid.New().String(), filepath.Ext(image.FileHeader.Filename))
		if err := c.SaveUploadedFile(image.FileHeader, "src"+newFileName); err != nil {
			return err
		} else {
			image.File = newFileName
			image.OriginalName = image.FileHeader.Filename
		}
	}

	var imageOld *models.GalleryImage
	imageRepo := models.NewGalleryImageRepository()
	if image.ID != 0 {
		imageOld, _ = imageRepo.GetByID(image.ID)
	}

	if imageOld == nil || imageOld.ID == 0 {
		err := imageRepo.Create(image)
		if err != nil {
			return err
		}
	} else {
		err := imageRepo.Update(image)
		if err != nil {
			return err
		}
	}

	return nil
}
