package repository

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/logger"
	common "github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"gorm.io/gorm"
)

type GalleryRepository interface {
	Create(gallery *models.Gallery) error
	GetByID(id uint) (*models.Gallery, error)
	Update(gallery *models.Gallery) error
	Delete(id uint) error
	GetAll() ([]models.Gallery, error)
	GetAllIds() ([]uint, error)
	GetAllForResource(contracts.Resource) ([]*models.Gallery, error)
}

type galleryRepository struct {
	*common.Paginate
	db *gorm.DB
}

func NewGalleryRepository() GalleryRepository {
	return &galleryRepository{db: db.GetDB()}
}

func (r *galleryRepository) Create(gallery *models.Gallery) error {
	return r.db.Create(gallery).Error
}

func (r *galleryRepository) GetByID(id uint) (*models.Gallery, error) {
	var gallery models.Gallery
	if err := r.db.First(&gallery, id).Error; err != nil {
		return nil, err
	}
	return &gallery, nil
}

func (r *galleryRepository) Update(gallery *models.Gallery) error {
	return r.db.Save(gallery).Error
}

func (r *galleryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Gallery{}, id).Error
}

func (r *galleryRepository) GetAll() ([]models.Gallery, error) {
	var galleries []models.Gallery
	if err := r.db.Find(&galleries).Error; err != nil {
		return nil, err
	}
	return galleries, nil
}

func (r *galleryRepository) GetAllForResource(c contracts.Resource) ([]*models.Gallery, error) {
	var galleries []*models.Gallery
	err := r.db.
		Joins("inner join gallery_has_resources as r on galleries.id = r.gallery_id").
		Where("r.resource_id = ?", c.GetID()).
		Where("r.resource = ?", c.GetResource()).
		Model(&models.Gallery{}).
		Preload("GalleryImage").
		Find(&galleries).Error
	return galleries, err
}

func (r *galleryRepository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.Gallery{}).Pluck("id", &ids).Error; err != nil {
		logger.New().Error(err)
	}
	return ids, nil
}
