package repository

import (
	"github.com/axlle-com/blog/pkg/common/db"
	common "github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/template/models"
	"gorm.io/gorm"
	"log"
)

type TemplateRepository interface {
	Create(template *models.Template) error
	GetByID(id uint) (*models.Template, error)
	Update(template *models.Template) error
	Delete(id uint) error
	GetAll() ([]*models.Template, error)
	GetAllIds() ([]uint, error)
}

type repository struct {
	*common.Paginate
	db *gorm.DB
}

func NewRepo() TemplateRepository {
	return &repository{db: db.GetDB()}
}

func (r *repository) Create(template *models.Template) error {
	return r.db.Create(template).Error
}

func (r *repository) GetByID(id uint) (*models.Template, error) {
	var template models.Template
	if err := r.db.First(&template, id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *repository) Update(template *models.Template) error {
	return r.db.Save(template).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&models.Template{}, id).Error
}

func (r *repository) GetAll() ([]*models.Template, error) {
	var template []*models.Template
	if err := r.db.Find(&template).Error; err != nil {
		return nil, err
	}
	return template, nil
}

func (r *repository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.Template{}).Pluck("id", &ids).Error; err != nil {
		log.Printf("Failed to fetch IDs from the database: %v\n", err)
	}
	return ids, nil
}
