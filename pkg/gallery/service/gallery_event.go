package service

import (
	"github.com/axlle-com/blog/pkg/gallery/models"
	"github.com/axlle-com/blog/pkg/gallery/repository"
)

type GalleryEvent struct {
	imageService *ImageService
	resourceRepo repository.GalleryResourceRepository
}

func NewGalleryEvent(
	image *ImageService,
	resource repository.GalleryResourceRepository,
) *GalleryEvent {
	return &GalleryEvent{
		imageService: image,
		resourceRepo: resource,
	}
}

func (e *GalleryEvent) DeletingGallery(g *models.Gallery) (err error) {
	has, _ := e.resourceRepo.GetByGalleryID(g.ID)
	if has != nil {
		if err = e.resourceRepo.Delete(g.ID); err != nil {
			return err
		}
	}

	err = e.imageService.DeleteImages(g.Images)
	if err != nil {
		return err
	}
	return nil
}

func (e *GalleryEvent) DeletedGallery(g *models.Gallery) (err error) {
	return err
}
