package repository

import (
	"github.com/axlle-com/blog/pkg/common/db"
	common "github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/user/models"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(post *models.Role) error
	GetByID(id uint) (*models.Role, error)
	GetByName(name string) (*models.Role, error)
	Update(post *models.Role) error
	Delete(id uint) error
	GetAll() ([]models.Role, error)
	GetPaginate(page, pageSize int) ([]models.Role, error)
}

type roleRepository struct {
	*common.Paginate
	db *gorm.DB
}

func NewRoleRepository() RoleRepository {
	return &roleRepository{db: db.GetDB()}
}

func (r *roleRepository) Create(post *models.Role) error {
	return r.db.Create(post).Error
}

func (r *roleRepository) GetByID(id uint) (*models.Role, error) {
	var model models.Role
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *roleRepository) GetByName(name string) (*models.Role, error) {
	var model models.Role
	if err := r.db.Where("name = ?", name).First(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *roleRepository) Update(model *models.Role) error {
	return r.db.Save(model).Error
}

func (r *roleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}

func (r *roleRepository) GetAll() ([]models.Role, error) {
	var m []models.Role
	if err := r.db.Find(&m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r *roleRepository) GetPaginate(page, pageSize int) ([]models.Role, error) {
	var m []models.Role

	err := r.db.Scopes(r.SetPaginate(page, pageSize)).Find(&m).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}
