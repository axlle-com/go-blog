package service

import (
	"errors"
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/gallery/models"
)

func CreateGallery(g *models.Gallery) (*models.Gallery, error) {
	repo := models.GalleryRepo()

	if err := repo.Create(g); err != nil {
		return nil, err
	}

	err := galleryImageUpdate(g)
	return g, err
}

func UpdateGallery(g *models.Gallery) (*models.Gallery, error) {
	repo := models.GalleryRepo()

	if err := repo.Update(g); err != nil {
		return nil, err
	}

	err := galleryImageUpdate(g)
	return g, err
}

func DeleteGalleries(g []*models.Gallery) (err error) {
	var ids []uint
	var resGalleries []*models.Gallery
	isErr := false
	for _, gallery := range g {
		if err = DeletingGallery(gallery); err == nil {
			ids = append(ids, gallery.ID)
			resGalleries = append(resGalleries, gallery)
		} else {
			isErr = true
			logger.Error(err)
		}
	}

	if isErr {
		return errors.New("Ошибки при удалении галлерей")
	}

	if len(ids) > 0 {
		if err = models.GalleryRepo().DeleteByIDs(ids); err == nil {
			for _, gallery := range resGalleries {
				if err = DeletedGallery(gallery); err != nil {
					isErr = true
					logger.Error(err)
				}
			}

			if isErr {
				return errors.New("Ошибки в удалении галлерей")
			}
		}
	}
	return
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
			image, e := SaveImage(item)
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
