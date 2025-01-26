package repository

import (
	"github.com/axlle-com/blog/pkg/app/db"
	"github.com/axlle-com/blog/pkg/app/logger"
	app "github.com/axlle-com/blog/pkg/app/models"
	"github.com/axlle-com/blog/pkg/post/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	WithTx(tx *gorm.DB) CategoryRepository
	Create(postCategory *models.PostCategory) error
	GetByID(id uint) (*models.PostCategory, error)
	Update(postCategory *models.PostCategory) error
	Delete(id uint) error
	GetAll() ([]*models.PostCategory, error)
	GetAllIds() ([]uint, error)
	WithPaginate(page, pageSize int) ([]*models.PostCategory, error)
}

type categoryRepository struct {
	db *gorm.DB
	*app.Paginate
}

func NewCategoryRepo() CategoryRepository {
	r := &categoryRepository{db: db.GetDB()}
	return r
}

func (r *categoryRepository) WithTx(tx *gorm.DB) CategoryRepository {
	newR := &categoryRepository{db: tx}
	return newR
}

func (r *categoryRepository) Create(postCategory *models.PostCategory) error {
	return r.db.Create(postCategory).Error
}

func (r *categoryRepository) GetByID(id uint) (*models.PostCategory, error) {
	var postCategory models.PostCategory
	if err := r.db.First(&postCategory, id).Error; err != nil {
		return nil, err
	}
	return &postCategory, nil
}

func (r *categoryRepository) Update(postCategory *models.PostCategory) error {
	return r.db.Save(postCategory).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.PostCategory{}, id).Error
}

func (r *categoryRepository) GetAll() ([]*models.PostCategory, error) {
	var postCategories []*models.PostCategory
	if err := r.db.Find(&postCategories).Error; err != nil {
		return nil, err
	}
	return postCategories, nil
}

func (r *categoryRepository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.PostCategory{}).Pluck("id", &ids).Error; err != nil {
		logger.Error(err)
	}
	return ids, nil
}

func (r *categoryRepository) WithPaginate(page, pageSize int) ([]*models.PostCategory, error) {
	var models []*models.PostCategory

	err := r.db.Model(models).Scopes(r.SetPaginate(page, pageSize)).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}
