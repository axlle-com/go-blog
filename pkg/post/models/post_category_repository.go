package models

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/logger"
	common "github.com/axlle-com/blog/pkg/common/models"
)

type CategoryRepository interface {
	Create(postCategory *PostCategory) error
	GetByID(id uint) (*PostCategory, error)
	Update(postCategory *PostCategory) error
	Delete(id uint) error
	GetAll() ([]*PostCategory, error)
	GetAllIds() ([]uint, error)
	GetPaginate(page, pageSize int) ([]*PostCategory, error)
	Transaction()
	Rollback()
	Commit()
}

type categoryRepository struct {
	*common.Repo
	*common.Paginate
}

func CategoryRepo() CategoryRepository {
	r := &categoryRepository{Repo: &common.Repo{}}
	r.SetConnection(db.GetDB())
	return r
}

func (r *categoryRepository) Create(postCategory *PostCategory) error {
	return r.Connection().Create(postCategory).Error
}

func (r *categoryRepository) GetByID(id uint) (*PostCategory, error) {
	var postCategory PostCategory
	if err := r.Connection().First(&postCategory, id).Error; err != nil {
		return nil, err
	}
	return &postCategory, nil
}

func (r *categoryRepository) Update(postCategory *PostCategory) error {
	return r.Connection().Save(postCategory).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.Connection().Delete(&PostCategory{}, id).Error
}

func (r *categoryRepository) GetAll() ([]*PostCategory, error) {
	var postCategories []*PostCategory
	if err := r.Connection().Find(&postCategories).Error; err != nil {
		return nil, err
	}
	return postCategories, nil
}

func (r *categoryRepository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.Connection().Model(&PostCategory{}).Pluck("id", &ids).Error; err != nil {
		logger.Error(err)
	}
	return ids, nil
}

func (r *categoryRepository) GetPaginate(page, pageSize int) ([]*PostCategory, error) {
	var models []*PostCategory

	err := r.Connection().Model(models).Scopes(r.SetPaginate(page, pageSize)).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}
