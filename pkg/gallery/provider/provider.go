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
	err = service.DeleteForResource(c)
	if err != nil {
		return err
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
