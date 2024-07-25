package repository

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/models"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	Create(post *models.Permission) error
	GetByID(id uint) (*models.Permission, error)
	Update(post *models.Permission) error
	Delete(id uint) error
	GetAll() ([]models.Permission, error)
	GetPaginate(page, pageSize int) ([]models.Permission, error)
}

type repository struct {
	*models.Paginate
	db *gorm.DB
}

func NewPermissionRepository() PermissionRepository {
	return &repository{db: db.GetDB()}
}

func (r *repository) Create(post *models.Permission) error {
	return r.db.Create(post).Error
}

func (r *repository) GetByID(id uint) (*models.Permission, error) {
	var model models.Permission
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *repository) Update(model *models.Permission) error {
	return r.db.Save(model).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&models.Permission{}, id).Error
}

func (r *repository) GetAll() ([]models.Permission, error) {
	var models []models.Permission
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *repository) GetPaginate(page, pageSize int) ([]models.Permission, error) {
	var models []models.Permission

	err := r.db.Scopes(r.SetPaginate(page, pageSize)).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}
