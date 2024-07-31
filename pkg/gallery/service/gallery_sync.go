package service

import (
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"github.com/axlle-com/blog/pkg/gallery/repository"
)

func Attach(r contracts.Resource, g *models.Gallery) error {
	hasRepo, err := repository.
		NewGalleryHasResourceRepository().
		GetByResourceAndID(r.GetID(), r.GetResource(), g.ID)
	if err != nil || hasRepo == nil {
		err = repository.
			NewGalleryHasResourceRepository().
			Create(
				&models.GalleryHasResource{
					ResourceID: r.GetID(),
					Resource:   r.GetResource(),
					GalleryID:  g.ID,
				},
			)
	}
	return err
}

func Detach(r contracts.Resource, g *models.Gallery) error {
	hasRepo, err := repository.
		NewGalleryHasResourceRepository().
		GetByResourceAndID(r.GetID(), r.GetResource(), g.ID)
	if err != nil || hasRepo == nil {
		err = repository.
			NewGalleryHasResourceRepository().
			Create(
				&models.GalleryHasResource{
					ResourceID: r.GetID(),
					Resource:   r.GetResource(),
					GalleryID:  g.ID,
				},
			)
	}
	return err
}
