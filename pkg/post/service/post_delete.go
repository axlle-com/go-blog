package service

import (
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	"github.com/axlle-com/blog/pkg/post/models"
)

func PostDelete(post *models.Post) error {
	err := gallery.Provider().DeleteForResource(post)
	if err != nil {
		return err
	}
	if err := models.PostRepo().Delete(post); err != nil {
		return err
	}
	return nil
}
