package models

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/logger"
	common "github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"gorm.io/gorm"
)

type GalleryRepository interface {
	Create(gallery *Gallery) error
	GetByID(id uint) (*Gallery, error)
	GetByIDs(ids []uint) ([]*Gallery, error)
	Update(gallery *Gallery) error
	DeleteByID(id uint) error
	Delete(*Gallery) error
	DeleteByIDs(ids []uint) (err error)
	GetAll() ([]*Gallery, error)
	GetAllIds() ([]uint, error)
	GetForResource(contracts.Resource) ([]*Gallery, error)
	WithImages() GalleryRepository
}

type galleryRepository struct {
	*common.Paginate
	db         *gorm.DB
	withImages bool
}

func GalleryRepo() GalleryRepository {
	return &galleryRepository{db: db.GetDB()}
}

func (r *galleryRepository) WithImages() GalleryRepository {
	r.withImages = true
	return r
}

func (r *galleryRepository) Create(gallery *Gallery) error {
	return r.db.Omit("Images").Create(gallery).Error
}

func (r *galleryRepository) GetByID(id uint) (*Gallery, error) {
	var gallery Gallery
	if err := r.db.First(&gallery, id).Error; err != nil {
		return nil, err
	}
	return &gallery, nil
}

func (r *galleryRepository) GetByIDs(ids []uint) ([]*Gallery, error) {
	var galleries []*Gallery
	query := r.db.Where("id IN ?", ids)

	if r.withImages {
		query.Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort ASC")
		})
	}

	if err := query.Find(&galleries).Error; err != nil {
		return nil, err
	}
	return galleries, nil
}

func (r *galleryRepository) Update(gallery *Gallery) error {
	return r.db.Select("Title", "Description", "Sort", "Image", "URL").Save(gallery).Error
}

func (r *galleryRepository) DeleteByID(id uint) error {
	return r.db.Delete(Gallery{}, id).Error
}

func (r *galleryRepository) Delete(g *Gallery) (err error) {
	return r.db.Delete(g, g.ID).Error
}

func (r *galleryRepository) DeleteByIDs(ids []uint) (err error) {
	return r.db.Where("id IN ?", ids).Delete(&Gallery{}).Error
}

func (r *galleryRepository) GetAll() ([]*Gallery, error) {
	var galleries []*Gallery
	if err := r.db.Find(&galleries).Error; err != nil {
		return nil, err
	}
	return galleries, nil
}

func (r *galleryRepository) GetForResource(c contracts.Resource) ([]*Gallery, error) {
	var galleries []*Gallery
	query := r.db.
		Joins("inner join gallery_has_resources as r on galleries.id = r.gallery_id").
		Where("r.resource_id = ?", c.GetID()).
		Where("r.resource = ?", c.GetResource()).
		Model(&Gallery{})

	if r.withImages {
		query.Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort ASC")
		})
	}

	err := query.Find(&galleries).Error
	return galleries, err
}

func (r *galleryRepository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&Gallery{}).Pluck("id", &ids).Error; err != nil {
		logger.Error(err)
	}
	return ids, nil
}
