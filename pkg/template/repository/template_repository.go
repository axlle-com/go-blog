package repository

import (
	"github.com/axlle-com/blog/pkg/common/db"
	common "github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/template/models"
	"log"
)

type TemplateRepository interface {
	Create(template *models.Template) error
	GetByID(id uint) (*models.Template, error)
	Update(template *models.Template) error
	Delete(id uint) error
	GetAll() ([]*models.Template, error)
	GetAllIds() ([]uint, error)
	Transaction()
	Rollback()
	Commit()
}

type repository struct {
	*common.Repo
	*common.Paginate
}

func NewRepo() TemplateRepository {
	r := &repository{Repo: &common.Repo{}}
	r.SetConnection(db.GetDB())
	return r
}

func (r *repository) Create(template *models.Template) error {
	return r.Connection().Create(template).Error
}

func (r *repository) GetByID(id uint) (*models.Template, error) {
	var template models.Template
	if err := r.Connection().First(&template, id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *repository) Update(template *models.Template) error {
	return r.Connection().Save(template).Error
}

func (r *repository) Delete(id uint) error {
	return r.Connection().Delete(&models.Template{}, id).Error
}

func (r *repository) GetAll() ([]*models.Template, error) {
	var template []*models.Template
	if err := r.Connection().Find(&template).Error; err != nil {
		return nil, err
	}
	return template, nil
}

func (r *repository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.Connection().Model(&models.Template{}).Pluck("id", &ids).Error; err != nil {
		log.Printf("Failed to fetch IDs from the database: %v\n", err)
	}
	return ids, nil
}
