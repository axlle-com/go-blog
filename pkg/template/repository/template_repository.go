package repository

import (
	"fmt"
	"log"

	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/template/http/request"
	"github.com/axlle-com/blog/pkg/template/models"
	"gorm.io/gorm"
)

type TemplateRepository interface {
	WithTx(tx *gorm.DB) TemplateRepository
	Create(template *models.Template) error
	GetByID(id uint) (*models.Template, error)
	GetByIDs(ids []uint) ([]*models.Template, error)
	GetByNames(names []string) ([]*models.Template, error)
	GetByName(name string) (*models.Template, error)
	FindByFilter(filter *models.TemplateFilter) (*models.Template, error)
	Update(template *models.Template) error
	Delete(*models.Template) error
	DeleteByIDs(ids []uint) (err error)
	GetAll() ([]*models.Template, error)
	GetAllIds() ([]uint, error)
	WithPaginate(p contract.Paginator, filter *request.TemplateFilter) ([]*models.Template, error)
	Filter(filter *models.TemplateFilter) ([]*models.Template, error)
}

type repository struct {
	db *gorm.DB
	*app.Paginate
}

func NewTemplateRepo(db *gorm.DB) TemplateRepository {
	r := &repository{db: db}
	return r
}

func (r *repository) WithTx(tx *gorm.DB) TemplateRepository {
	newR := &repository{db: tx}
	return newR
}

func (r *repository) Create(template *models.Template) error {
	template.Creating()
	return r.db.Create(template).Error
}

func (r *repository) GetByID(id uint) (*models.Template, error) {
	var template models.Template
	if err := r.db.First(&template, id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *repository) GetByIDs(ids []uint) ([]*models.Template, error) {
	var templates []*models.Template
	if err := r.db.Where("id IN (?)", ids).Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *repository) GetByNames(names []string) ([]*models.Template, error) {
	if len(names) == 0 {
		return []*models.Template{}, nil
	}
	var templates []*models.Template
	if err := r.db.Where("name IN ?", names).Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *repository) GetByName(name string) (*models.Template, error) {
	var template models.Template
	query := r.db.Where("name = ?", name)

	if err := query.First(&template).Error; err != nil {
		return nil, err
	}

	return &template, nil
}

func (r *repository) FindByFilter(filter *models.TemplateFilter) (*models.Template, error) {
	var template models.Template
	query := r.db.Model(&template)

	// Применяем фильтры напрямую из полей структуры, если они не nil
	if filter.ID != nil {
		query = query.Where("id = ?", *filter.ID)
	}

	if filter.Name != nil {
		query = query.Where("name = ?", *filter.Name)
	}

	if filter.ResourceName != nil {
		if *filter.ResourceName != "" {
			query = query.Where("resource_name = ?", *filter.ResourceName)
		} else {
			query = query.Where("resource_name IS NULL OR resource_name = ''")
		}
	}

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	if filter.Title != nil {
		query = query.Where("title = ?", *filter.Title)
	}

	if err := query.First(&template).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *repository) Update(template *models.Template) error {
	template.Updating()
	return r.db.Select(
		"UserID",
		"Title",
		"IsMain",
		"Name",
		"ResourceName",
		"HTML",
		"JS",
		"CSS",
	).Save(template).Error
}

func (r *repository) Delete(template *models.Template) error {
	return r.db.Delete(&models.Template{}, template.ID).Error
}

func (r *repository) DeleteByIDs(ids []uint) (err error) {
	return r.db.Where("id IN ?", ids).Delete(&models.Template{}).Error
}

func (r *repository) GetAll() ([]*models.Template, error) {
	var template []*models.Template
	if err := r.db.Order("id ASC").Find(&template).Error; err != nil {
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

func (r *repository) WithPaginate(paginator contract.Paginator, filter *request.TemplateFilter) ([]*models.Template, error) {
	var templates []*models.Template
	var total int64

	template := models.Template{}
	table := template.GetTable()

	query := r.db.Model(&template)

	// TODO WHERE IN; LIKE
	for col, val := range filter.GetMap() {
		if col == "title" {
			query = query.Where(fmt.Sprintf("%s.%v ilike ?", table, col), fmt.Sprintf("%%%v%%", val))
			continue
		}
		query = query.Where(fmt.Sprintf("%s.%v = ?", table, col), val)
	}

	query.Count(&total)

	err := query.Scopes(r.SetPaginate(paginator.GetPage(), paginator.GetPageSize())).
		Order(fmt.Sprintf("%s.id ASC", template.GetTable())).
		Find(&templates).Error
	if err != nil {
		return nil, err
	}

	paginator.SetTotal(int(total))
	return templates, nil
}

func (r *repository) Filter(filter *models.TemplateFilter) ([]*models.Template, error) {
	var templates []*models.Template
	template := models.Template{}

	query := r.db.Model(&template)

	// Применяем фильтры, если они не nil
	if filter.ID != nil {
		query = query.Where("id = ?", *filter.ID)
	}

	if filter.Name != nil {
		query = query.Where("name = ?", *filter.Name)
	}

	if filter.ResourceName != nil && *filter.ResourceName != "" {
		query = query.Where("resource_name = ?", *filter.ResourceName)
	}

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	if filter.Title != nil {
		query = query.Where("title = ?", *filter.Title)
	}

	err := query.Order(fmt.Sprintf("%s.id ASC", template.GetTable())).Find(&templates).Error
	if err != nil {
		return nil, err
	}

	return templates, nil
}
