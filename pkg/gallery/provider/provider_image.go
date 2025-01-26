package provider

import (
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	"github.com/axlle-com/blog/pkg/gallery/repository"
)

type ImageProvider interface {
	GetForGallery(id uint) []contracts.Image
	GetAll() []contracts.Image
}

func NewImageProvider(
	image repository.GalleryImageRepository,
) ImageProvider {
	return &imageProvider{
		image: image,
	}
}

type imageProvider struct {
	image repository.GalleryImageRepository
}

func (p *imageProvider) GetForGallery(id uint) []contracts.Image {
	var collection []contracts.Image
	images, err := p.image.GetForGallery(id)
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
	images, err := p.image.GetAll()
	if err == nil {
		for _, image := range images {
			collection = append(collection, image)
		}
		return collection
	}
	return nil
}
