package service

import (
	"errors"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"github.com/axlle-com/blog/pkg/common/service"
	"github.com/axlle-com/blog/pkg/gallery/models"
)

func SaveGallery(g any) (contracts.Gallery, error) {
	gal := service.LoadStruct(&models.Gallery{}, g).(*models.Gallery)
	repo := models.GalleryRepo()

	if gal.ID == 0 {
		if err := repo.Create(gal); err != nil {
			return nil, err
		}
	} else {
		if err := repo.Update(gal); err != nil {
			return nil, err
		}
	}

	var err error
	if len(gal.Images) > 0 {
		slice := make([]*models.Image, len(gal.Images), len(gal.Images))
		var eSlice []error
		for _, image := range gal.Images {
			image.GalleryID = gal.ID
			i, e := SaveImage(image)
			if e != nil {
				eSlice = append(eSlice, e)
				continue
			}
			slice = append(slice, i)
		}
		if len(eSlice) > 0 {
			err = errors.New("Были ошибки при сохранении изображения")
		}
		gal.Images = slice
	}

	return gal, err
}
