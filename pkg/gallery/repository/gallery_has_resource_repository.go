package repository

import (
	"github.com/axlle-com/blog/pkg/common/db"
	common "github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"gorm.io/gorm"
)

type GalleryHasResourceRepository interface {
	GetByResourceAndID(id uint, resource string, galleryID uint) (*models.GalleryHasResource, error)
	Create(*models.GalleryHasResource) error
}

type galleryHasResourceRepository struct {
	*common.Paginate
	db *gorm.DB
}

func NewGalleryHasResourceRepository() GalleryHasResourceRepository {
	return &galleryHasResourceRepository{db: db.GetDB()}
}

func (r *galleryHasResourceRepository) Create(galleryHasResource *models.GalleryHasResource) error {
	return r.db.Create(galleryHasResource).Error
}

func (r *galleryHasResourceRepository) GetByResourceAndID(id uint, resource string, galleryID uint) (*models.GalleryHasResource, error) {
	var galleryHasResource models.GalleryHasResource
	if err := r.db.
		Where("resource_id = ?", id).
		Where("resource = ?", resource).
		Where("gallery_id = ?", galleryID).
		First(&galleryHasResource).Error; err != nil {
		return nil, err
	}
	return &galleryHasResource, nil
}
