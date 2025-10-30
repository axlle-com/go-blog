package provider

import (
	"github.com/axlle-com/blog/app/models/contract"
	apppPovider "github.com/axlle-com/blog/app/models/provider"
	"github.com/axlle-com/blog/pkg/gallery/repository"
)

func NewImageProvider(
	image repository.GalleryImageRepository,
) apppPovider.ImageProvider {
	return &imageProvider{
		image: image,
	}
}

type imageProvider struct {
	image repository.GalleryImageRepository
}

func (p *imageProvider) GetForGallery(id uint) []contract.Image {
	var collection []contract.Image
	images, err := p.image.GetForGallery(id)
	if err == nil {
		for _, image := range images {
			collection = append(collection, image)
		}
		return collection
	}
	return nil
}

func (p *imageProvider) GetAll() []contract.Image {
	var collection []contract.Image
	images, err := p.image.GetAll()
	if err == nil {
		for _, image := range images {
			collection = append(collection, image)
		}
		return collection
	}
	return nil
}
