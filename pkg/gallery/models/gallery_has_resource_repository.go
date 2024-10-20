package models

import (
	"errors"
	"github.com/axlle-com/blog/pkg/common/db"
	common "github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"gorm.io/gorm"
)

type GalleryResourceRepository interface {
	GetByResourceAndID(id uint, resource string, galleryID uint) (*GalleryHasResource, error)
	GetForResource(contracts.Resource) ([]*GalleryHasResource, error)
	GetByGalleryID(uint) (*GalleryHasResource, error)
	Create(*GalleryHasResource) error
	Delete(uint) error
	DetachResource(contracts.Resource) error
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

func (r *galleryResourceRepository) GetForResource(c contracts.Resource) ([]*GalleryHasResource, error) {
	var galleryHasResource []*GalleryHasResource
	err := r.db.
		Where("resource = ?", c.GetResource()).
		Where("resource_id = ?", c.GetID()).
		Find(&galleryHasResource).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return galleryHasResource, nil
}

func (r *galleryResourceRepository) DetachResource(c contracts.Resource) error {
	err := r.db.
		Where("resource = ?", c.GetResource()).
		Where("resource_id = ?", c.GetID()).
		Delete(&GalleryHasResource{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

func (r *galleryResourceRepository) GetByGalleryID(id uint) (*GalleryHasResource, error) {
	var galleryHasResource GalleryHasResource
	if err := r.db.
		Where("gallery_id = ?", id).
		First(&galleryHasResource).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &galleryHasResource, nil
}

func (r *galleryResourceRepository) Delete(id uint) error {
	return r.db.Where("gallery_id = ?", id).Delete(&GalleryHasResource{}).Error
}
