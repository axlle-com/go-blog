package service

import (
	"fmt"
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/file"
	"github.com/axlle-com/blog/pkg/gallery/models"
)

func _SaveImage(image *models.Image) error {
	if image.FileHeader != nil {
		newFileName := fmt.Sprintf("gallery/%d", image.GalleryID)
		path, err := file.SaveUploadedFile(image.FileHeader, newFileName)
		if err != nil {
			logger.Error(err)
			return err
		} else {
			image.File = path
			image.OriginalName = image.FileHeader.Filename
		}
	}

	var imageOld *models.Image
	imageRepo := models.ImageRepo()
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
