package service

import (
	"errors"
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/gallery/models"
)

func SaveImage(image *models.Image) (*models.Image, error) {
	repo := models.ImageRepo()

	if image.ID == 0 {
		if err := repo.Create(image); err != nil {
			return nil, err
		}
	} else {
		if err := repo.Update(image); err != nil {
			return nil, err
		}
	}

	return image, nil
}

func DeleteImages(is []*models.Image) (err error) {
	var ids []uint
	var resImages []*models.Image
	var isErr bool
	for _, im := range is {
		if err := DeletingImage(im); err == nil {
			ids = append(ids, im.ID)
			resImages = append(resImages, im)
		} else {
			isErr = true
			logger.Error(err)
		}
	}

	if isErr {
		return errors.New("Ошибки при удалении изображений")
	}

	if len(ids) > 0 {
		if err = models.ImageRepo().DeleteByIDs(ids); err == nil {
			for _, im := range resImages {
				if err := DeletedImage(im); err != nil {
					logger.Error(err)
					isErr = true
				}
			}

			if isErr {
				return errors.New("Ошибки при удалении изображений")
			}
		}
	}
	return
}

func DeleteImage(im *models.Image) (err error) {
	if err = DeletingImage(im); err != nil {
		return
	}

	if err = models.ImageRepo().Delete(im); err == nil {
		if err = DeletedImage(im); err != nil {
			return
		}
		return
	}
	return
}
