package provider

import (
	"github.com/axlle-com/blog/pkg/app/logger"
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	app "github.com/axlle-com/blog/pkg/app/service"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"github.com/axlle-com/blog/pkg/gallery/repository"
	"github.com/axlle-com/blog/pkg/gallery/service"
)

type GalleryProvider interface {
	GetForResource(contracts.Resource) []contracts.Gallery
	GetAll() []contracts.Gallery
	SaveFromForm(g any, resource contracts.Resource) (contracts.Gallery, error)
	DeleteForResource(contracts.Resource) error
}

func NewProvider(
	gallery repository.GalleryRepository,
	service *service.GalleryService,
) GalleryProvider {
	return &provider{
		gallery: gallery,
		service: service,
	}
}

type provider struct {
	gallery repository.GalleryRepository
	service *service.GalleryService
}

func (p *provider) GetForResource(resource contracts.Resource) []contracts.Gallery {
	galleries, err := p.gallery.WithImages().GetForResource(resource)
	collection := make([]contracts.Gallery, 0, len(galleries))
	if err == nil {
		for _, gallery := range galleries {
			collection = append(collection, gallery)
		}
		return collection
	}
	logger.Error(err)
	return nil
}

func (p *provider) DeleteForResource(resource contracts.Resource) (err error) {
	err = p.service.DeleteForResource(resource)
	if err != nil {
		return err
	}

	return nil
}

func (p *provider) GetAll() []contracts.Gallery {
	var collection []contracts.Gallery
	galleries, err := p.gallery.GetAll()
	if err == nil {
		for _, gallery := range galleries {
			collection = append(collection, gallery)
		}
		return collection
	}
	logger.Error(err)
	return nil
}

func (p *provider) SaveFromForm(g any, resource contracts.Resource) (gallery contracts.Gallery, err error) {
	gal := app.LoadStruct(&models.Gallery{}, g).(*models.Gallery)
	if gal.ID == 0 {
		gallery, err = p.service.CreateGallery(gal)
	} else {
		gallery, err = p.service.UpdateGallery(gal)
	}

	if err != nil {
		return nil, err
	}

	err = p.service.Attach(resource, gallery)
	if err != nil {
		return nil, err
	}

	return gallery, nil
}
