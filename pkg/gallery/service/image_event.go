package service

import (
	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/pkg/gallery/models"
)

type ImageEvent struct {
	galleryEvent *GalleryEvent
	api          *api.Api
}

func NewImageEvent(
	api *api.Api,
) *ImageEvent {
	return &ImageEvent{
		api: api,
	}
}

func (e *ImageEvent) DeletingImage(im *models.Image) (err error) {
	return
}

func (e *ImageEvent) SetGalleryEvent(galleryEvent *GalleryEvent) {
	if galleryEvent == nil {
		return
	}

	e.galleryEvent = galleryEvent
}

func (e *ImageEvent) DeletedImage(im *models.Image) (err error) {
	err = e.api.File.DeleteFile(im.File)
	if err != nil {
		return err
	}

	if e.galleryEvent == nil {
		return
	}
	e.galleryEvent.UpdateTrigger([]uint{im.GalleryID})

	return nil
}
