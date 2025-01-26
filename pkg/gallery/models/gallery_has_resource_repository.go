package models

import (
	"errors"
	"github.com/axlle-com/blog/pkg/common/db"
	common "github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"gorm.io/gorm"
)

type GalleryResourceRepository interface {
	WithTx(tx *gorm.DB) GalleryResourceRepository
	GetByParams(resourceID uint, resource string, galleryID uint) (*GalleryHasResource, error)
	DeleteByParams(resourceID uint, resource string, galleryID uint) error
	GetForResource(contracts.Resource) ([]*GalleryHasResource, error)
	GetByGalleryID(uint) (*GalleryHasResource, error)
	GetByResource(c contracts.Resource) ([]*GalleryHasResource, error)
	Create(*GalleryHasResource) error
	Delete(uint) error
	DetachResource(contracts.Resource) error
}

type galleryResourceRepository struct {
	db *gorm.DB
	*common.Paginate
}

func ResourceRepo() GalleryResourceRepository {
	r := &galleryResourceRepository{db: db.GetDB()}
	return r
}

func (r *galleryResourceRepository) WithTx(tx *gorm.DB) GalleryResourceRepository {
	newR := &galleryResourceRepository{db: tx}
	return newR
}

func (r *galleryResourceRepository) Create(galleryHasResource *GalleryHasResource) error {
	return r.db.Create(galleryHasResource).Error
}

func (r *galleryResourceRepository) GetByParams(resourceID uint, resource string, galleryID uint) (*GalleryHasResource, error) {
	var galleryHasResource GalleryHasResource
	if err := r.db.
		Where("resource_id = ?", resourceID).
		Where("resource = ?", resource).
		Where("gallery_id = ?", galleryID).
		First(&galleryHasResource).Error; err != nil {
		return nil, err
	}
	return &galleryHasResource, nil
}

func (r *galleryResourceRepository) DeleteByParams(resourceID uint, resource string, galleryID uint) error {
	err := r.db.
		Where("resource_id = ?", resourceID).
		Where("resource = ?", resource).
		Where("gallery_id = ?", galleryID).
		Delete(&GalleryHasResource{}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return err
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

func (r *galleryResourceRepository) GetByResource(c contracts.Resource) ([]*GalleryHasResource, error) {
	var galleryHasResource []*GalleryHasResource
	err := r.db.
		Where("resource_id = ? AND resource = ?", c.GetID(), c.GetResource()).
		Or("gallery_id IN (?)",
			r.db.Model(&GalleryHasResource{}).
				Select("gallery_id").
				Where("resource_id = ? AND resource = ?", c.GetID(), c.GetResource()),
		).
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
	err := r.db.Where("gallery_id = ?", id).Delete(&GalleryHasResource{}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}
