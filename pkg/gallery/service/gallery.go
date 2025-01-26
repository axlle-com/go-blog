package service

import (
	"errors"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
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

func DeleteForResource(c contracts.Resource) (err error) {
	rRepo := models.ResourceRepo()
	resource, err := rRepo.GetByResource(c)
	if err != nil {
		return err
	}

	all := make(map[uint]*models.GalleryHasResource)
	only := make(map[uint]*models.GalleryHasResource)
	detach := make(map[uint]*models.GalleryHasResource)
	var galleryIDs []uint
	if resource == nil {
		return nil
	}

	for _, r := range resource {
		if r.ResourceID != c.GetID() && r.Resource != c.GetResource() {
			all[r.GalleryID] = r
		} else {
			only[r.GalleryID] = r
		}
	}

	for id, _ := range only {
		if _, ok := all[id]; ok {
			detach[id] = all[id]
		} else {
			galleryIDs = append(galleryIDs, id)
		}
	}

	if len(detach) > 0 { // TODO need test
		for _, r := range detach {
			err = rRepo.DeleteByParams(r.ResourceID, r.Resource, r.GalleryID)
			if err != nil {
				return err
			}
		}
	}

	if len(galleryIDs) > 0 {
		galleries, err := models.GalleryRepo().WithImages().GetByIDs(galleryIDs)
		if err != nil {
			return err
		}
		err = DeleteGalleries(galleries)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteGalleries(galleries []*models.Gallery) (err error) {
	var ids []uint
	for _, gallery := range galleries {
		if err = DeletingGallery(gallery); err != nil {
			return err
		}
		ids = append(ids, gallery.ID)
	}

	if len(ids) > 0 {
		if err = models.GalleryRepo().DeleteByIDs(ids); err == nil {
			for _, gallery := range galleries {
				if err = DeletedGallery(gallery); err != nil {
					return err
				}
			}
			return nil
		}
	}
	return err
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
