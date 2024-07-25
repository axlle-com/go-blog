package repository

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/models"
	"gorm.io/gorm"
	"log"
)

type Repository interface {
	Create(postCategory *models.PostCategory) error
	GetByID(id uint) (*models.PostCategory, error)
	Update(postCategory *models.PostCategory) error
	Delete(id uint) error
	GetAll() ([]models.PostCategory, error)
	GetAllIds() ([]uint, error)
	GetPaginate(page, pageSize int) ([]models.PostCategory, error)
}

type repository struct {
	*models.Paginate
	db *gorm.DB
}

func NewRepository() Repository {
	return &repository{db: db.GetDB()}
}

func (r *repository) Create(postCategory *models.PostCategory) error {
	return r.db.Create(postCategory).Error
}

func (r *repository) GetByID(id uint) (*models.PostCategory, error) {
	var postCategory models.PostCategory
	if err := r.db.First(&postCategory, id).Error; err != nil {
		return nil, err
	}
	return &postCategory, nil
}

func (r *repository) Update(postCategory *models.PostCategory) error {
	return r.db.Save(postCategory).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&models.PostCategory{}, id).Error
}

func (r *repository) GetAll() ([]models.PostCategory, error) {
	var postCategories []models.PostCategory
	if err := r.db.Find(&postCategories).Error; err != nil {
		return nil, err
	}
	return postCategories, nil
}

func (r *repository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.PostCategory{}).Pluck("id", &ids).Error; err != nil {
		log.Println("Failed to fetch IDs from the database: %v", err)
	}
	return ids, nil
}

func (r *repository) GetPaginate(page, pageSize int) ([]models.PostCategory, error) {
	var models []models.PostCategory

	err := r.db.Model(models).Scopes(r.SetPaginate(page, pageSize)).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}
