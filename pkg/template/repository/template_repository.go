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

type templateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository() TemplateRepository {
	return &templateRepository{db: db.GetDB()}
}

func (r *templateRepository) CreateTemplate(template *models.Template) error {
	return r.db.Create(template).Error
}

func (r *templateRepository) GetTemplateByID(id uint) (*models.Template, error) {
	var template models.Template
	if err := r.db.First(&template, id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *templateRepository) UpdateTemplate(template *models.Template) error {
	return r.db.Save(template).Error
}

func (r *templateRepository) DeleteTemplate(id uint) error {
	return r.db.Delete(&models.Template{}, id).Error
}

func (r *templateRepository) GetAllTemplates() ([]models.Template, error) {
	var template []models.Template
	if err := r.db.Find(&template).Error; err != nil {
		return nil, err
	}
	return template, nil
}

func (r *templateRepository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.Template{}).Pluck("id", &ids).Error; err != nil {
		log.Println("Failed to fetch IDs from the database: %v", err)
	}
	return ids, nil
}
