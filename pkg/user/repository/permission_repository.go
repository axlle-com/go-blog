package repository

import (
	"github.com/axlle-com/blog/app/db"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/pkg/user/models"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	WithTx(tx *gorm.DB) PermissionRepository
	Create(post *models.Permission) error
	GetByID(id uint) (*models.Permission, error)
	Update(post *models.Permission) error
	Delete(id uint) error
	GetAll() ([]models.Permission, error)
	WithPaginate(page, pageSize int) ([]models.Permission, error)
}

type permissionRepository struct {
	db *gorm.DB
	*app.Paginate
}

func NewPermissionRepo() PermissionRepository {
	return &permissionRepository{db: db.GetDB()}
}

func (r *permissionRepository) WithTx(tx *gorm.DB) PermissionRepository {
	newR := &permissionRepository{db: tx}
	return newR
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

func (r *permissionRepository) WithPaginate(page, pageSize int) ([]models.Permission, error) {
	var m []models.Permission

	err := r.db.Scopes(r.SetPaginate(page, pageSize)).Find(&m).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}
