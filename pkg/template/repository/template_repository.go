package repository

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/models"
	"gorm.io/gorm"
	"log"
)

type TemplateRepository interface {
	CreateTemplate(template *models.Template) error
	GetTemplateByID(id uint) (*models.Template, error)
	UpdateTemplate(template *models.Template) error
	DeleteTemplate(id uint) error
	GetAllTemplates() ([]models.Template, error)
	GetAllIds() ([]uint, error)
}

type repository struct {
	*models.Paginate
	db *gorm.DB
}

func NewRepository() TemplateRepository {
	return &repository{db: db.GetDB()}
}

func (r *repository) CreateTemplate(template *models.Template) error {
	return r.db.Create(template).Error
}

func (r *repository) GetTemplateByID(id uint) (*models.Template, error) {
	var template models.Template
	if err := r.db.First(&template, id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *repository) UpdateTemplate(template *models.Template) error {
	return r.db.Save(template).Error
}

func (r *repository) DeleteTemplate(id uint) error {
	return r.db.Delete(&models.Template{}, id).Error
}

func (r *repository) GetAllTemplates() ([]models.Template, error) {
	var template []models.Template
	if err := r.db.Find(&template).Error; err != nil {
		return nil, err
	}
	return template, nil
}

func (r *repository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.Template{}).Pluck("id", &ids).Error; err != nil {
		log.Println("Failed to fetch IDs from the database: %v", err)
	}
	return ids, nil
}
