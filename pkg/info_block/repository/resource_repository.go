package repository

import (
	"errors"

	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InfoBlockHasResourceRepository interface {
	WithTx(tx *gorm.DB) InfoBlockHasResourceRepository
	Create(*models.InfoBlockHasResource) error
	Update(infoBlockHasResource *models.InfoBlockHasResource) error
	FindByParams(resourceUUID uuid.UUID, infoBlockID uint) (*models.InfoBlockHasResource, error)
	FindByID(id uint) (*models.InfoBlockHasResource, error)
	DeleteByID(uint) error
	DeleteByResourceUUID(uuid.UUID) error
	DeleteByInfoBlockID(infoBlockID uint) error
}

type infoBlockResource struct {
	db *gorm.DB
}

func NewResourceRepo(db *gorm.DB) InfoBlockHasResourceRepository {
	r := &infoBlockResource{db: db}
	return r
}

func (r *infoBlockResource) WithTx(tx *gorm.DB) InfoBlockHasResourceRepository {
	return &infoBlockResource{db: tx}
}

func (r *infoBlockResource) Create(infoBlockHasResource *models.InfoBlockHasResource) error {
	return r.db.Create(infoBlockHasResource).Error
}

func (r *infoBlockResource) FindByID(id uint) (*models.InfoBlockHasResource, error) {
	var infoBlockHasResource models.InfoBlockHasResource
	if err := r.db.
		Where("id = ?", id).
		First(&infoBlockHasResource).Error; err != nil {
		return nil, err
	}
	return &infoBlockHasResource, nil
}

func (r *infoBlockResource) FindByParams(resourceUUID uuid.UUID, infoBlockID uint) (*models.InfoBlockHasResource, error) {
	var infoBlockHasResource models.InfoBlockHasResource
	if err := r.db.
		Where("resource_uuid = ?", resourceUUID).
		Where("info_block_id = ?", infoBlockID).
		First(&infoBlockHasResource).Error; err != nil {
		return nil, err
	}
	return &infoBlockHasResource, nil
}

func (r *infoBlockResource) DeleteByID(id uint) error {
	err := r.db.
		Where("id = ?", id).
		Delete(&models.InfoBlockHasResource{}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

func (r *infoBlockResource) DeleteByResourceUUID(resourceUUID uuid.UUID) error {
	err := r.db.
		Where("resource_uuid = ?", resourceUUID).
		Delete(&models.InfoBlockHasResource{}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

func (r *infoBlockResource) DeleteByInfoBlockID(infoBlockID uint) error {
	err := r.db.
		Where("info_block_id = ?", infoBlockID).
		Delete(&models.InfoBlockHasResource{}).Error

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
