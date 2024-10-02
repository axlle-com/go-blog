package models

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/logger"
	common "github.com/axlle-com/blog/pkg/common/models"
	"gorm.io/gorm"
)

type GalleryImageRepository interface {
	Create(image *GalleryImage) error
	GetByID(id uint) (*GalleryImage, error)
	Update(image *GalleryImage) error
	Delete(id uint) error
	GetAll() ([]GalleryImage, error)
	GetAllIds() ([]uint, error)
	CountForGallery(id uint) int64
}

type galleryImageRepository struct {
	*common.Paginate
	db *gorm.DB
}

func NewGalleryImageRepository() GalleryImageRepository {
	return &galleryImageRepository{db: db.GetDB()}
}

func (r *galleryImageRepository) Create(image *GalleryImage) error {
	return r.db.Create(image).Error
}

func (r *galleryImageRepository) GetByID(id uint) (*GalleryImage, error) {
	var image GalleryImage
	if err := r.db.First(&image, id).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *galleryImageRepository) Update(image *GalleryImage) error {
	return r.db.Select("GalleryID", "Title", "Description", "Sort").Save(image).Error
}

// Delete TODO транзакции
func (r *galleryImageRepository) Delete(id uint) error {
	g, err := r.GetByID(id)
	if err == nil {
		if err = r.db.Delete(&GalleryImage{}, id).Error; err == nil {
			return g.Deleted()
		}
	}
	return err
}

func (r *galleryImageRepository) GetAll() ([]GalleryImage, error) {
	var images []GalleryImage
	if err := r.db.Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func (r *galleryImageRepository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&GalleryImage{}).Pluck("id", &ids).Error; err != nil {
		logger.Error(err)
	}
	return ids, nil
}

func (r *galleryImageRepository) CountForGallery(id uint) int64 {
	var count int64
	result := r.db.Model(&GalleryImage{}).Where("gallery_id = ?", id).Count(&count)
	if result.Error != nil {
		logger.Error(result.Error)
	}
	return count
}
