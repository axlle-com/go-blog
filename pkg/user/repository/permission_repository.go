package repository

import (
	"github.com/axlle-com/blog/pkg/common/db"
	common "github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/user/models"
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

type permissionRepository struct {
	*common.Paginate
	db *gorm.DB
}

func NewPermissionRepository() PermissionRepository {
	return &permissionRepository{db: db.GetDB()}
}

func (r *permissionRepository) Create(post *models.Permission) error {
	return r.db.Create(post).Error
}

func (r *permissionRepository) GetByID(id uint) (*models.Permission, error) {
	var model models.Permission
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *permissionRepository) Update(model *models.Permission) error {
	return r.db.Save(model).Error
}

func (r *permissionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Permission{}, id).Error
}

func (r *permissionRepository) GetAll() ([]models.Permission, error) {
	var m []models.Permission
	if err := r.db.Find(&m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r *permissionRepository) GetPaginate(page, pageSize int) ([]models.Permission, error) {
	var m []models.Permission

	err := r.db.Scopes(r.SetPaginate(page, pageSize)).Find(&m).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}
