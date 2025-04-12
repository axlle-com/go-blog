package repository

import (
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"gorm.io/gorm"
)

type GalleryImageRepository interface {
	WithTx(tx *gorm.DB) GalleryImageRepository
	Create(image *models.Image) error
	GetByID(id uint) (*models.Image, error)
	Update(image *models.Image) error
	Delete(image *models.Image) error
	DeleteByIDs(ids []uint) (err error)
	GetAll() ([]*models.Image, error)
	GetAllIds() ([]uint, error)
	CountForGallery(id uint) int64
	GetForGallery(id uint) ([]*models.Image, error)
}

type galleryImageRepository struct {
	db *gorm.DB
	*app.Paginate
}

func NewImageRepo() GalleryImageRepository {
	r := &galleryImageRepository{db: db.GetDB()}
	return r
}

func (r *galleryImageRepository) WithTx(tx *gorm.DB) GalleryImageRepository {
	newR := &galleryImageRepository{db: tx}
	return newR
}

func (r *galleryImageRepository) Create(image *models.Image) error {
	return r.db.Create(image).Error
}

func (r *galleryImageRepository) GetByID(id uint) (*models.Image, error) {
	var image models.Image
	if err := r.db.First(&image, id).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *galleryImageRepository) Update(image *models.Image) error {
	return r.db.Select("GalleryID", "Title", "Description", "Sort").Save(image).Error
}

func (r *galleryImageRepository) Delete(image *models.Image) (err error) {
	return r.db.Delete(&models.Image{}, image.ID).Error
}

func (r *galleryImageRepository) DeleteByIDs(ids []uint) (err error) {
	return r.db.Where("id IN ?", ids).Delete(&models.Image{}).Error
}

func (r *galleryImageRepository) GetAll() ([]*models.Image, error) {
	var images []*models.Image
	if err := r.db.Order("id ASC").Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func (r *galleryImageRepository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.Image{}).Pluck("id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *galleryImageRepository) CountForGallery(id uint) int64 {
	var count int64
	result := r.db.Model(&models.Image{}).Where("gallery_id = ?", id).Count(&count)
	if result.Error != nil {
		logger.Error(result.Error)
	}
	return count
}

func (r *galleryImageRepository) GetForGallery(id uint) ([]*models.Image, error) {
	var images []*models.Image
	if err := r.db.Model(&models.Image{}).Where("gallery_id = ?", id).Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}
