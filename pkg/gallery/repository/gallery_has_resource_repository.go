package repository

import (
	"errors"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GalleryResourceRepository interface {
	WithTx(tx *gorm.DB) GalleryResourceRepository
	GetByParams(resourceUUID uuid.UUID, galleryID uint) (*models.GalleryHasResource, error)
	DeleteByParams(resourceUUID uuid.UUID, galleryID uint) error
	GetForResource(contracts.Resource) ([]*models.GalleryHasResource, error)
	GetByGalleryID(uint) (*models.GalleryHasResource, error)
	GetByResource(c contracts.Resource) ([]*models.GalleryHasResource, error)
	Create(*models.GalleryHasResource) error
	Delete(uint) error
	DetachResource(contracts.Resource) error
}

type galleryResourceRepository struct {
	db *gorm.DB
	*app.Paginate
}

func NewResourceRepo(db contracts.DB) GalleryResourceRepository {
	r := &galleryResourceRepository{db: db.GORM()}
	return r
}

func (r *galleryResourceRepository) WithTx(tx *gorm.DB) GalleryResourceRepository {
	newR := &galleryResourceRepository{db: tx}
	return newR
}

func (r *galleryResourceRepository) Create(galleryHasResource *models.GalleryHasResource) error {
	return r.db.Create(galleryHasResource).Error
}

func (r *galleryResourceRepository) GetByParams(resourceUUID uuid.UUID, galleryID uint) (*models.GalleryHasResource, error) {
	var galleryHasResource models.GalleryHasResource
	if err := r.db.
		Where("resource_uuid = ?", resourceUUID).
		Where("gallery_id = ?", galleryID).
		First(&galleryHasResource).Error; err != nil {
		return nil, err
	}
	return &galleryHasResource, nil
}

func (r *galleryResourceRepository) DeleteByParams(resourceUUID uuid.UUID, galleryID uint) error {
	err := r.db.
		Where("resource_uuid = ?", resourceUUID).
		Where("gallery_id = ?", galleryID).
		Delete(&models.GalleryHasResource{}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return err
}

func (r *galleryResourceRepository) GetForResource(resource contracts.Resource) ([]*models.GalleryHasResource, error) {
	var galleryHasResource []*models.GalleryHasResource
	err := r.db.
		Where("resource_uuid = ?", resource.GetUUID()).
		Find(&galleryHasResource).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return galleryHasResource, nil
}

func (r *galleryResourceRepository) GetByResource(resource contracts.Resource) ([]*models.GalleryHasResource, error) {
	var galleryHasResource []*models.GalleryHasResource
	err := r.db.
		Where("resource_uuid = ?", resource.GetUUID()).
		Or("gallery_id IN (?)",
			r.db.Model(&models.GalleryHasResource{}).
				Select("gallery_id").
				Where("resource_uuid = ?", resource.GetUUID()),
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

func (r *galleryResourceRepository) DetachResource(resource contracts.Resource) error {
	err := r.db.
		Where("resource_uuid = ?", resource.GetUUID()).
		Delete(&models.GalleryHasResource{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

func (r *galleryResourceRepository) GetByGalleryID(id uint) (*models.GalleryHasResource, error) {
	var galleryHasResource models.GalleryHasResource
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
	err := r.db.Where("gallery_id = ?", id).Delete(&models.GalleryHasResource{}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}
