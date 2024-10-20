package service

import (
	"github.com/axlle-com/blog/pkg/file"
	"github.com/axlle-com/blog/pkg/gallery/models"
)

func DeletingImage(im *models.Image) (err error) {
	return
}

func DeletedImage(im *models.Image) (err error) {
	err = file.DeleteFile(im.File)
	if err != nil {
		return err
	}
	return
}
