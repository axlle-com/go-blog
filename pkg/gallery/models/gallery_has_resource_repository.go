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
	DeleteByResourceAndID(id uint, resource string, galleryID uint) error
	GetForResource(contracts.Resource) ([]*GalleryHasResource, error)
	GetByGalleryID(uint) (*GalleryHasResource, error)
	GetGalleriesByResource(c contracts.Resource) ([]*GalleryHasResource, error)
	Create(*GalleryHasResource) error
	Delete(uint) error
	DetachResource(contracts.Resource) error
	Transaction()
	Rollback()
	Commit()
}

type galleryResourceRepository struct {
	*common.Repo
	*common.Paginate
}

func ResourceRepo() GalleryResourceRepository {
	r := &galleryResourceRepository{Repo: &common.Repo{}}
	r.SetConnection(db.GetDB())
	return r
}

func (r *galleryResourceRepository) Create(galleryHasResource *GalleryHasResource) error {
	return r.Connection().Create(galleryHasResource).Error
}

func (r *galleryResourceRepository) GetByResourceAndID(id uint, resource string, galleryID uint) (*GalleryHasResource, error) {
	var galleryHasResource GalleryHasResource
	if err := r.Connection().
		Where("resource_id = ?", id).
		Where("resource = ?", resource).
		Where("gallery_id = ?", galleryID).
		First(&galleryHasResource).Error; err != nil {
		return nil, err
	}
	return &galleryHasResource, nil
}

func (r *galleryResourceRepository) DeleteByResourceAndID(id uint, resource string, galleryID uint) error {
	err := r.Connection().
		Where("resource_id = ?", id).
		Where("resource = ?", resource).
		Where("gallery_id = ?", galleryID).
		Delete(&GalleryHasResource{}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	if err != nil {
		r.Rollback()
	}

	return err
}

func (r *galleryResourceRepository) GetForResource(c contracts.Resource) ([]*GalleryHasResource, error) {
	var galleryHasResource []*GalleryHasResource
	err := r.Connection().
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

func (r *galleryResourceRepository) GetGalleriesByResource(c contracts.Resource) ([]*GalleryHasResource, error) {
	var galleryHasResource []*GalleryHasResource
	err := r.Connection().
		Where("resource_id = ? AND resource = ?", c.GetID(), c.GetResource()).
		Or("gallery_id IN (?)",
			r.Connection().Model(&GalleryHasResource{}).
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
	err := r.Connection().
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
	if err := r.Connection().
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
	err := r.Connection().Where("gallery_id = ?", id).Delete(&GalleryHasResource{}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}
