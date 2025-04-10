package repository

import (
	"errors"
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InfoBlockHasResourceRepository interface {
	WithTx(tx *gorm.DB) InfoBlockHasResourceRepository
	GetByParams(resourceUUID uuid.UUID, infoBlockID uint) (*models.InfoBlockHasResource, error)
	DeleteByParams(resourceUUID uuid.UUID, infoBlockID uint) error
	GetByID(id uint) (*models.InfoBlockHasResource, error)
	GetForResource(contracts.Resource) ([]*models.InfoBlockHasResource, error)
	GetByGalleryID(uint) (*models.InfoBlockHasResource, error)
	GetByResource(c contracts.Resource) ([]*models.InfoBlockHasResource, error)
	Create(*models.InfoBlockHasResource) error
	Delete(uint) error
	DetachResource(contracts.Resource) error
	Update(infoBlockHasResource *models.InfoBlockHasResource) error
}

type infoBlockResource struct {
	db *gorm.DB
}

func NewResourceRepo() InfoBlockHasResourceRepository {
	r := &infoBlockResource{db: db.GetDB()}
	return r
}

func (r *infoBlockResource) WithTx(tx *gorm.DB) InfoBlockHasResourceRepository {
	return &infoBlockResource{db: tx}
}

func (r *infoBlockResource) Create(infoBlockHasResource *models.InfoBlockHasResource) error {
	return r.db.Create(infoBlockHasResource).Error
}

func (r *infoBlockResource) GetByID(id uint) (*models.InfoBlockHasResource, error) {
	var infoBlockHasResource models.InfoBlockHasResource
	if err := r.db.
		Where("id = ?", id).
		First(&infoBlockHasResource).Error; err != nil {
		return nil, err
	}
	return &infoBlockHasResource, nil
}

func (r *infoBlockResource) GetByParams(resourceUUID uuid.UUID, infoBlockID uint) (*models.InfoBlockHasResource, error) {
	var infoBlockHasResource models.InfoBlockHasResource
	if err := r.db.
		Where("resource_uuid = ?", resourceUUID).
		Where("info_block_id = ?", infoBlockID).
		First(&infoBlockHasResource).Error; err != nil {
		return nil, err
	}
	return &infoBlockHasResource, nil
}

func (r *infoBlockResource) DeleteByParams(resourceUUID uuid.UUID, infoBlockID uint) error {
	err := r.db.
		Where("resource_uuid = ?", resourceUUID).
		Where("info_block_id = ?", infoBlockID).
		Delete(&models.InfoBlockHasResource{}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return err
}

func (r *infoBlockResource) GetForResource(resource contracts.Resource) ([]*models.InfoBlockHasResource, error) {
	var infoBlockHasResource []*models.InfoBlockHasResource
	err := r.db.
		Where("resource_uuid = ?", resource.GetUUID()).
		Find(&infoBlockHasResource).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return infoBlockHasResource, nil
}

func (r *infoBlockResource) GetByResource(resource contracts.Resource) ([]*models.InfoBlockHasResource, error) {
	var infoBlockHasResource []*models.InfoBlockHasResource
	err := r.db.
		Where("resource_uuid = ?", resource.GetUUID()).
		Or("info_block_id IN (?)",
			r.db.Model(&models.InfoBlockHasResource{}).
				Select("info_block_id").
				Where("resource_uuid = ?", resource.GetUUID()),
		).
		Find(&infoBlockHasResource).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return infoBlockHasResource, nil
}

func (r *infoBlockResource) DetachResource(resource contracts.Resource) error {
	err := r.db.
		Where("resource_uuid = ?", resource.GetUUID()).
		Delete(&models.InfoBlockHasResource{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

func (r *infoBlockResource) GetByGalleryID(id uint) (*models.InfoBlockHasResource, error) {
	var infoBlockHasResource models.InfoBlockHasResource
	if err := r.db.
		Where("info_block_id = ?", id).
		First(&infoBlockHasResource).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &infoBlockHasResource, nil
}

func (r *infoBlockResource) Delete(id uint) error {
	err := r.db.Where("info_block_id = ?", id).Delete(&models.InfoBlockHasResource{}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

func (r *infoBlockResource) Update(infoBlockHasResource *models.InfoBlockHasResource) error {
	return r.db.Select(
		"Position",
		"Sort",
	).Save(infoBlockHasResource).Error
}
