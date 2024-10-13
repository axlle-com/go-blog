package models

import (
	"github.com/axlle-com/blog/pkg/common/db"
	common "github.com/axlle-com/blog/pkg/common/models"
	"gorm.io/gorm"
)

type GalleryResourceRepository interface {
	GetByResourceAndID(id uint, resource string, galleryID uint) (*GalleryHasResource, error)
	GetByID(uint) (*GalleryHasResource, error)
	Create(*GalleryHasResource) error
	Delete(uint) error
}

type galleryResourceRepository struct {
	*common.Paginate
	db *gorm.DB
}

func ResourceRepo() GalleryResourceRepository {
	return &galleryResourceRepository{db: db.GetDB()}
}

func (r *galleryResourceRepository) Create(galleryHasResource *GalleryHasResource) error {
	return r.db.Create(galleryHasResource).Error
}

func (r *galleryResourceRepository) GetByResourceAndID(id uint, resource string, galleryID uint) (*GalleryHasResource, error) {
	var galleryHasResource GalleryHasResource
	if err := r.db.
		Where("resource_id = ?", id).
		Where("resource = ?", resource).
		Where("gallery_id = ?", galleryID).
		First(&galleryHasResource).Error; err != nil {
		return nil, err
	}
	return &galleryHasResource, nil
}

func (r *galleryResourceRepository) GetByID(id uint) (*GalleryHasResource, error) {
	var galleryHasResource GalleryHasResource
	if err := r.db.
		Where("gallery_id = ?", id).
		First(&galleryHasResource).Error; err != nil {
		return nil, err
	}
	return &galleryHasResource, nil
}

func (r *galleryResourceRepository) Delete(id uint) error {
	return r.db.Where("gallery_id = ?", id).Delete(&GalleryHasResource{}).Error
}
