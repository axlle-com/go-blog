package models

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/logger"
	common "github.com/axlle-com/blog/pkg/common/models"
)

type GalleryImageRepository interface {
	Create(image *Image) error
	GetByID(id uint) (*Image, error)
	Update(image *Image) error
	Delete(image *Image) error
	DeleteByIDs(ids []uint) (err error)
	GetAll() ([]*Image, error)
	GetAllIds() ([]uint, error)
	CountForGallery(id uint) int64
	GetForGallery(id uint) ([]*Image, error)
	Transaction()
	Rollback()
	Commit()
}

type galleryImageRepository struct {
	*common.Repo
	*common.Paginate
}

func ImageRepo() GalleryImageRepository {
	r := &galleryImageRepository{Repo: &common.Repo{}}
	r.SetConnection(db.GetDB())
	return r
}

func (r *galleryImageRepository) Create(image *Image) error {
	return r.Connection().Create(image).Error
}

func (r *galleryImageRepository) GetByID(id uint) (*Image, error) {
	var image Image
	if err := r.Connection().First(&image, id).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *galleryImageRepository) Update(image *Image) error {
	return r.Connection().Select("GalleryID", "Title", "Description", "Sort").Save(image).Error
}

func (r *galleryImageRepository) Delete(image *Image) (err error) {
	return r.Connection().Delete(&Image{}, image.ID).Error
}

func (r *galleryImageRepository) DeleteByIDs(ids []uint) (err error) {
	return r.Connection().Where("id IN ?", ids).Delete(&Image{}).Error
}

func (r *galleryImageRepository) GetAll() ([]*Image, error) {
	var images []*Image
	if err := r.Connection().Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func (r *galleryImageRepository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.Connection().Model(&Image{}).Pluck("id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *galleryImageRepository) CountForGallery(id uint) int64 {
	var count int64
	result := r.Connection().Model(&Image{}).Where("gallery_id = ?", id).Count(&count)
	if result.Error != nil {
		logger.Error(result.Error)
	}
	return count
}

func (r *galleryImageRepository) GetForGallery(id uint) ([]*Image, error) {
	var images []*Image
	if err := r.Connection().Model(&Image{}).Where("gallery_id = ?", id).Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}
