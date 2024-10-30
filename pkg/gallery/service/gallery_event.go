package service

import (
	"github.com/axlle-com/blog/pkg/gallery/models"
)

func DeletingGallery(g *models.Gallery) (err error) {
	repo := models.ResourceRepo()
	has, _ := repo.GetByGalleryID(g.ID)
	if has != nil {
		if err = repo.Delete(g.ID); err != nil {
			return err
		}
	}

	err = DeleteImages(g.Images)
	if err != nil {
		return err
	}
	return nil
}

func DeletedGallery(g *models.Gallery) (err error) {
	return err
}
