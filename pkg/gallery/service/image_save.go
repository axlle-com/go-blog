package service

import (
	"github.com/axlle-com/blog/pkg/common/service"
	"github.com/axlle-com/blog/pkg/gallery/models"
)

func SaveImage(i any) (*models.Image, error) {
	image := service.LoadStruct(&models.Image{}, i).(*models.Image)
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
