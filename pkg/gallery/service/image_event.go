package service

import (
	"github.com/axlle-com/blog/pkg/file/provider"
	"github.com/axlle-com/blog/pkg/gallery/models"
)

type ImageEvent struct {
	fileProvider provider.FileProvider
}

func NewImageEvent(
	file provider.FileProvider,
) *ImageEvent {
	return &ImageEvent{
		fileProvider: file,
	}
}

func (e *ImageEvent) DeletingImage(im *models.Image) (err error) {
	return
}

func (e *ImageEvent) DeletedImage(im *models.Image) (err error) {
	err = e.fileProvider.DeleteFile(im.File)
	if err != nil {
		return err
	}
	return
}
