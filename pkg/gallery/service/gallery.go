package service

import (
	"errors"
	"github.com/axlle-com/blog/pkg/gallery/models"
)

func GalleryCreate(g *models.Gallery) (*models.Gallery, error) {
	repo := models.GalleryRepo()

	if err := repo.Create(g); err != nil {
		return nil, err
	}

	err := galleryImageUpdate(g)
	return g, err
}

func GalleryUpdate(g *models.Gallery) (*models.Gallery, error) {
	repo := models.GalleryRepo()

	if err := repo.Update(g); err != nil {
		return nil, err
	}

	err := galleryImageUpdate(g)
	return g, err
}

func galleryImageUpdate(g *models.Gallery) error {
	var err error
	if len(g.Images) > 0 {
		slice := make([]*models.Image, len(g.Images), len(g.Images))
		var eSlice []error
		for idx, item := range g.Images {
			if item == nil {
				continue
			}
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

	return err
}
