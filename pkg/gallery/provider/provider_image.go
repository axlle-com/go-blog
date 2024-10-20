package provider

import (
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"github.com/axlle-com/blog/pkg/gallery/models"
)

type Image interface {
	GetForGallery(id uint) []contracts.Image
	GetAll() []contracts.Image
}

func ImageProvider() Image {
	return &imageProvider{}
}

type imageProvider struct {
}

func (p *imageProvider) GetForGallery(id uint) []contracts.Image {
	var collection []contracts.Image
	images, err := models.
		ImageRepo().
		GetForGallery(id)
	if err == nil {
		for _, image := range images {
			collection = append(collection, image)
		}
		return collection
	}
	return nil
}

func (p *imageProvider) GetAll() []contracts.Image {
	var collection []contracts.Image
	images, err := models.
		ImageRepo().
		GetAll()
	if err == nil {
		for _, image := range images {
			collection = append(collection, image)
		}
		return collection
	}
	return nil
}
