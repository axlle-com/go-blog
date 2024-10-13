package service

import (
	"errors"
	"github.com/axlle-com/blog/pkg/gallery/models"
)

func GallerySave(g *models.Gallery) (*models.Gallery, error) {
	repo := models.GalleryRepo()

	if g.ID == 0 {
		if err := repo.Create(g); err != nil {
			return nil, err
		}
	} else {
		if err := repo.Update(g); err != nil {
			return nil, err
		}
	}

	var err error
	if len(g.Images) > 0 {
		slice := make([]*models.Image, len(g.Images), len(g.Images))
		var eSlice []error
		for idx, item := range g.Images {
			item.GalleryID = g.ID
			image, e := ImageSave(item)
			if e != nil {
				eSlice = append(eSlice, e)
				continue
			}
			slice[idx] = image
		}
		if len(eSlice) > 0 {
			err = errors.New("Были ошибки при сохранении изображения")
		}
		g.Images = slice
	}

	return g, err
}
