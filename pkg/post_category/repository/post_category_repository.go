package repository

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/models"
	"gorm.io/gorm"
	"log"
)

type PostCategoryRepository interface {
	CreatePostCategory(postCategory *models.PostCategory) error
	GetPostCategoryByID(id uint) (*models.PostCategory, error)
	UpdatePostCategory(postCategory *models.PostCategory) error
	DeletePostCategory(id uint) error
	GetAllPostCategories() ([]models.PostCategory, error)
	GetAllIds() ([]uint, error)
}

type postCategoryRepository struct {
	db *gorm.DB
}

func NewPostCategoryRepository() PostCategoryRepository {
	return &postCategoryRepository{db: db.GetDB()}
}

func (r *postCategoryRepository) CreatePostCategory(postCategory *models.PostCategory) error {
	return r.db.Create(postCategory).Error
}

func (r *postCategoryRepository) GetPostCategoryByID(id uint) (*models.PostCategory, error) {
	var postCategory models.PostCategory
	if err := r.db.First(&postCategory, id).Error; err != nil {
		return nil, err
	}
	return &postCategory, nil
}

func (r *postCategoryRepository) UpdatePostCategory(postCategory *models.PostCategory) error {
	return r.db.Save(postCategory).Error
}

func (r *postCategoryRepository) DeletePostCategory(id uint) error {
	return r.db.Delete(&models.PostCategory{}, id).Error
}

func (r *postCategoryRepository) GetAllPostCategories() ([]models.PostCategory, error) {
	var postCategories []models.PostCategory
	if err := r.db.Find(&postCategories).Error; err != nil {
		return nil, err
	}
	return postCategories, nil
}

func (r *postCategoryRepository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.PostCategory{}).Pluck("id", &ids).Error; err != nil {
		log.Println("Failed to fetch IDs from the database: %v", err)
	}
	return ids, nil
}
