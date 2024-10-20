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

func (p *provider) DeleteForResource(c contracts.Resource) error {
	rRepo := models.ResourceRepo()
	resource, err := rRepo.GetForResource(c)
	if err != nil {
		return err
	}

	if resource != nil {
		if err := rRepo.DetachResource(c); err != nil {
			return err
		}
	}

	var galleryIDs []uint
	for _, g := range resource {
		rsc, err := rRepo.GetByGalleryID(g.GalleryID) // TODO
		if err != nil || rsc != nil {
			continue
		}
		galleryIDs = append(galleryIDs, g.GalleryID)
	}

	if len(galleryIDs) > 0 {
		galleries, err := models.GalleryRepo().WithImages().GetByIDs(galleryIDs)
		if err != nil {
			return err
		}
		err = service.DeleteGalleries(galleries)
		if err != nil {
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
