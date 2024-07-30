package repository

import (
	"github.com/axlle-com/blog/pkg/common/db"
	common "github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/gallery/models"
	"gorm.io/gorm"
	"log"
)

type GalleryImageRepository interface {
	Create(gallery *models.GalleryImage) error
	GetByID(id uint) (*models.GalleryImage, error)
	Update(gallery *models.GalleryImage) error
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

func (r *galleryImageRepository) Create(gallery *models.GalleryImage) error {
	return r.db.Create(gallery).Error
}

func (r *galleryImageRepository) GetByID(id uint) (*models.GalleryImage, error) {
	var gallery models.GalleryImage
	if err := r.db.First(&gallery, id).Error; err != nil {
		return nil, err
	}
	return &gallery, nil
}

func (r *galleryImageRepository) GetByEmail(email string) (*models.GalleryImage, error) {
	var gallery models.GalleryImage
	if err := r.db.Where("email = ?", email).First(&gallery).Error; err != nil {
		return nil, err
	}
	r.db.Preload("Roles.Permissions").Preload("Permissions").Find(&gallery)
	return &gallery, nil
}

func (r *galleryImageRepository) Update(gallery *models.GalleryImage) error {
	return r.db.Save(gallery).Error
}

func (r *galleryImageRepository) Delete(id uint) error {
	return r.db.Delete(&models.GalleryImage{}, id).Error
}

func (r *galleryImageRepository) GetAll() ([]models.GalleryImage, error) {
	var galleries []models.GalleryImage
	if err := r.db.Find(&galleries).Error; err != nil {
		return nil, err
	}
	return galleries, nil
}

func (r *galleryImageRepository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.GalleryImage{}).Pluck("id", &ids).Error; err != nil {
		log.Println("Failed to fetch IDs from the database: %v", err)
	}
	return ids, nil
}
