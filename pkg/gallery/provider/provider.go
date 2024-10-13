package provider

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	common "github.com/axlle-com/blog/pkg/common/service"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"github.com/axlle-com/blog/pkg/gallery/service"
)

type Gallery interface {
	GetAllForResource(contracts.Resource) []contracts.Gallery
	SaveFromForm(g any) (contracts.Gallery, error)
}

func Provider() Gallery {
	return &provider{}
}

type provider struct {
}

func (p *provider) GetAllForResource(c contracts.Resource) []contracts.Gallery {
	var collection []contracts.Gallery
	galleries, err := models.
		GalleryRepo().
		GetAllForResource(c)
	if err == nil {
		for _, gallery := range galleries {
			collection = append(collection, gallery)
		}
		return collection
	}
	logger.Error(err)
	return nil
}

func (p *provider) SaveFromForm(g any) (contracts.Gallery, error) {
	gal := common.LoadStruct(&models.Gallery{}, g).(*models.Gallery)
	save, err := service.GallerySave(gal)
	if err != nil {
		return nil, err
	}
	return save, nil
}
