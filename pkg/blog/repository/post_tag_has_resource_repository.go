package repository

import (
	"errors"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostTagResourceRepository interface {
	WithTx(tx *gorm.DB) PostTagResourceRepository
	GetByParams(resourceUUID uuid.UUID, postTagID uint) (*models.PostTagHasResource, error)
	DeleteByParams(resourceUUID uuid.UUID, postTagID uint) error
	DeleteByIDs(postTagIDs []uint) error
	GetForResource(resourceUUID uuid.UUID) ([]*models.PostTagHasResource, error)
	GetByPostTagID(uint) (*models.PostTagHasResource, error)
	GetByResource(c contracts.Resource) ([]*models.PostTagHasResource, error)
	Create(*models.PostTagHasResource) error
	Delete(uint) error
	DetachResource(contracts.Resource) error
}

type postTagResourceRepository struct {
	db *gorm.DB
	*app.Paginate
}

func NewResourceRepo(db contracts.DB) PostTagResourceRepository {
	r := &postTagResourceRepository{db: db.GORM()}
	return r
}

func (r *postTagResourceRepository) WithTx(tx *gorm.DB) PostTagResourceRepository {
	return &postTagResourceRepository{db: tx}
}

func (r *postTagResourceRepository) Create(postTagHasResource *models.PostTagHasResource) error {
	return r.db.Create(postTagHasResource).Error
}

func (r *postTagResourceRepository) GetByParams(resourceUUID uuid.UUID, postTagID uint) (*models.PostTagHasResource, error) {
	var postTagHasResource models.PostTagHasResource
	if err := r.db.
		Where("resource_uuid = ?", resourceUUID).
		Where("post_tag_id = ?", postTagID).
		First(&postTagHasResource).Error; err != nil {
		return nil, err
	}
	return &postTagHasResource, nil
}

func (r *postTagResourceRepository) DeleteByParams(resourceUUID uuid.UUID, postTagID uint) error {
	err := r.db.
		Where("resource_uuid = ?", resourceUUID).
		Where("post_tag_id = ?", postTagID).
		Delete(&models.PostTagHasResource{}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return err
}

func (r *postTagResourceRepository) DeleteByIDs(postTagIDs []uint) error {
	err := r.db.
		Where("post_tag_id IN ?", postTagIDs).
		Delete(&models.PostTagHasResource{}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return err
}

func (r *postTagResourceRepository) GetForResource(resourceUUID uuid.UUID) ([]*models.PostTagHasResource, error) {
	var postTagHasResource []*models.PostTagHasResource
	err := r.db.
		Where("resource_uuid = ?", resourceUUID).
		Find(&postTagHasResource).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return postTagHasResource, nil
}

func (r *postTagResourceRepository) GetByResource(resource contracts.Resource) ([]*models.PostTagHasResource, error) {
	var postTagHasResource []*models.PostTagHasResource
	err := r.db.
		Where("resource_uuid = ?", resource.GetUUID()).
		Or("post_tag_id IN (?)",
			r.db.Model(&models.PostTagHasResource{}).
				Select("post_tag_id").
				Where("resource_uuid = ?", resource.GetUUID()),
		).
		Find(&postTagHasResource).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return postTagHasResource, nil
}

func (r *postTagResourceRepository) DetachResource(resource contracts.Resource) error {
	err := r.db.
		Where("resource_uuid = ?", resource.GetUUID()).
		Delete(&models.PostTagHasResource{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

func (r *postTagResourceRepository) GetByPostTagID(id uint) (*models.PostTagHasResource, error) {
	var postTagHasResource models.PostTagHasResource
	if err := r.db.
		Where("post_tag_id = ?", id).
		First(&postTagHasResource).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &postTagHasResource, nil
}

func (r *postTagResourceRepository) Delete(id uint) error {
	err := r.db.Where("post_tag_id = ?", id).Delete(&models.PostTagHasResource{}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}
