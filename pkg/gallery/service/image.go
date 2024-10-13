package service

import (
	"github.com/axlle-com/blog/pkg/gallery/models"
)

func ImageSave(image *models.Image) (*models.Image, error) {
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
