package models

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/logger"
	common "github.com/axlle-com/blog/pkg/common/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(postCategory *PostCategory) error
	GetByID(id uint) (*PostCategory, error)
	Update(postCategory *PostCategory) error
	Delete(id uint) error
	GetAll() ([]*PostCategory, error)
	GetAllIds() ([]uint, error)
	GetPaginate(page, pageSize int) ([]*PostCategory, error)
}

type categoryRepository struct {
	*common.Paginate
	db *gorm.DB
}

func NewCategoryRepo() CategoryRepository {
	return &categoryRepository{db: db.GetDB()}
}

func (r *categoryRepository) Create(postCategory *PostCategory) error {
	return r.db.Create(postCategory).Error
}

func (r *categoryRepository) GetByID(id uint) (*PostCategory, error) {
	var postCategory PostCategory
	if err := r.db.First(&postCategory, id).Error; err != nil {
		return nil, err
	}
	return &postCategory, nil
}

func (r *categoryRepository) Update(postCategory *PostCategory) error {
	return r.db.Save(postCategory).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&PostCategory{}, id).Error
}

func (r *categoryRepository) GetAll() ([]*PostCategory, error) {
	var postCategories []*PostCategory
	if err := r.db.Find(&postCategories).Error; err != nil {
		return nil, err
	}
	return postCategories, nil
}

func (r *categoryRepository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&PostCategory{}).Pluck("id", &ids).Error; err != nil {
		logger.Error(err)
	}
	return ids, nil
}

func (r *categoryRepository) GetPaginate(page, pageSize int) ([]*PostCategory, error) {
	var models []*PostCategory

	err := r.db.Model(models).Scopes(r.SetPaginate(page, pageSize)).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}
