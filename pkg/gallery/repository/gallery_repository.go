package repository

import (
	"github.com/axlle-com/blog/pkg/app/db"
	"github.com/axlle-com/blog/pkg/app/logger"
	app "github.com/axlle-com/blog/pkg/app/models"
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"gorm.io/gorm"
)

type GalleryRepository interface {
	Create(gallery *models.Gallery) error
	GetByID(id uint) (*models.Gallery, error)
	GetByIDs(ids []uint) ([]*models.Gallery, error)
	Update(gallery *models.Gallery) error
	DeleteByID(id uint) error
	Delete(*models.Gallery) error
	DeleteByIDs(ids []uint) (err error)
	GetAll() ([]*models.Gallery, error)
	GetAllIds() ([]uint, error)
	GetForResource(contracts.Resource) ([]*models.Gallery, error)
	WithImages() GalleryRepository
	WithTx(tx *gorm.DB) GalleryRepository
}

type galleryRepository struct {
	db *gorm.DB
	*app.Paginate
	withImages bool
}

func NewGalleryRepo() GalleryRepository {
	r := &galleryRepository{db: db.GetDB()}
	return r
}

func (r *galleryRepository) WithTx(tx *gorm.DB) GalleryRepository {
	newR := &galleryRepository{db: tx}
	return newR
}

func (r *galleryRepository) WithImages() GalleryRepository {
	r.withImages = true
	return r
}

func (r *galleryRepository) Create(gallery *models.Gallery) error {
	return r.db.Omit("Images").Create(gallery).Error
}

func (r *galleryRepository) GetByID(id uint) (*models.Gallery, error) {
	var gallery models.Gallery
	if err := r.db.First(&gallery, id).Error; err != nil {
		return nil, err
	}
	return &gallery, nil
}

func (r *galleryRepository) GetByIDs(ids []uint) ([]*models.Gallery, error) {
	var galleries []*models.Gallery
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

func (r *galleryRepository) Update(gallery *models.Gallery) error {
	return r.db.Select("Title", "Description", "Sort", "Image", "URL").Save(gallery).Error
}

func (r *galleryRepository) DeleteByID(id uint) error {
	return r.db.Delete(models.Gallery{}, id).Error
}

func (r *galleryRepository) Delete(g *models.Gallery) (err error) {
	return r.db.Delete(g, g.ID).Error
}

func (r *galleryRepository) DeleteByIDs(ids []uint) (err error) {
	return r.db.Where("id IN ?", ids).Delete(&models.Gallery{}).Error
}

func (r *galleryRepository) GetAll() ([]*models.Gallery, error) {
	var galleries []*models.Gallery
	if err := r.db.Find(&galleries).Error; err != nil {
		return nil, err
	}
	return galleries, nil
}

func (r *galleryRepository) GetForResource(resource contracts.Resource) ([]*models.Gallery, error) {
	var galleries []*models.Gallery
	query := r.db.
		Joins("inner join gallery_has_resources as r on galleries.id = r.gallery_id").
		Where("r.resource_uuid = ?", resource.GetUUID()).
		Model(&models.Gallery{})

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
	if err := r.db.Model(&models.Gallery{}).Pluck("id", &ids).Error; err != nil {
		logger.Error(err)
	}
	return ids, nil
}
