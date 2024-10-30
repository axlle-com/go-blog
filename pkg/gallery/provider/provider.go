package provider

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	common "github.com/axlle-com/blog/pkg/common/service"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"github.com/axlle-com/blog/pkg/gallery/service"
)

type Gallery interface {
	GetForResource(contracts.Resource) []contracts.Gallery
	GetAll() []contracts.Gallery
	SaveFromForm(g any) (contracts.Gallery, error)
	DeleteForResource(contracts.Resource) error
}

func Provider() Gallery {
	return &provider{}
}

type provider struct {
}

func (p *provider) GetForResource(c contracts.Resource) []contracts.Gallery {
	var collection []contracts.Gallery
	galleries, err := models.
		GalleryRepo().
		WithImages().
		GetForResource(c)
	if err == nil {
		for _, gallery := range galleries {
			collection = append(collection, gallery)
		}
		return collection
	}
	logger.Error(err)
	return nil
}

func (p *provider) DeleteForResource(c contracts.Resource) (err error) {
	rRepo := models.ResourceRepo()
	resource, err := rRepo.GetGalleriesByResource(c)
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
		rRepo.Transaction()
		for _, r := range detach {
			err = rRepo.DeleteByResourceAndID(r.ResourceID, r.Resource, r.GalleryID)
			if err != nil {
				return err
			}
		}
	}

	if len(galleryIDs) > 0 {
		galleries, err := models.GalleryRepo().WithImages().GetByIDs(galleryIDs)
		if err != nil {
			rRepo.Rollback()
			return err
		}
		err = service.DeleteGalleries(galleries)
		if err != nil {
			rRepo.Rollback()
			return err
		}
	}
	return nil
}

func (p *provider) GetAll() []contracts.Gallery {
	var collection []contracts.Gallery
	galleries, err := models.
		GalleryRepo().
		GetAll()
	if err == nil {
		for _, gallery := range galleries {
			collection = append(collection, gallery)
		}
		return collection
	}
	logger.Error(err)
	return nil
}

func (p *provider) SaveFromForm(g any) (gallery contracts.Gallery, err error) {
	gal := common.LoadStruct(&models.Gallery{}, g).(*models.Gallery)
	if gal.ID == 0 {
		gallery, err = service.CreateGallery(gal)
	} else {
		gallery, err = service.UpdateGallery(gal)
	}

	if err != nil {
		return nil, err
	}
	return gallery, nil
}
