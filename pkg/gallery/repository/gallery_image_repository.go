package repository

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/logger"
	common "github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"gorm.io/gorm"
)

type GalleryImageRepository interface {
	Create(image *models.GalleryImage) error
	GetByID(id uint) (*models.GalleryImage, error)
	Update(image *models.GalleryImage) error
	Delete(id uint) error
	GetAll() ([]models.GalleryImage, error)
	GetAllIds() ([]uint, error)
}

type galleryImageRepository struct {
	*common.Paginate
	db *gorm.DB
}

func NewGalleryImageRepository() GalleryImageRepository {
	return &galleryImageRepository{db: db.GetDB()}
}

func (r *galleryImageRepository) Create(image *models.GalleryImage) error {
	return r.db.Create(image).Error
}

func (r *galleryImageRepository) GetByID(id uint) (*models.GalleryImage, error) {
	var image models.GalleryImage
	if err := r.db.First(&image, id).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *galleryImageRepository) GetByEmail(email string) (*models.GalleryImage, error) {
	var image models.GalleryImage
	if err := r.db.Where("email = ?", email).First(&image).Error; err != nil {
		return nil, err
	}
	r.db.Preload("Roles.Permissions").Preload("Permissions").Find(&image)
	return &image, nil
}

func (r *galleryImageRepository) Update(image *models.GalleryImage) error {
	return r.db.Select("GalleryID", "Title", "Description", "Sort").Save(image).Error
}

func (r *galleryImageRepository) Delete(id uint) error {
	return r.db.Delete(&models.GalleryImage{}, id).Error
}

func (r *galleryImageRepository) GetAll() ([]models.GalleryImage, error) {
	var images []models.GalleryImage
	if err := r.db.Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func (r *galleryImageRepository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.GalleryImage{}).Pluck("id", &ids).Error; err != nil {
		logger.New().Error(err)
	}
	return ids, nil
}
