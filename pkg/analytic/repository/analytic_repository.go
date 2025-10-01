package repository

import (
	"fmt"
	"log"

	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/analytic/models"
	"gorm.io/gorm"
)

type AnalyticRepository interface {
	WithTx(tx *gorm.DB) AnalyticRepository
	Create(analytic *models.Analytic) error
	GetByID(id uint) (*models.Analytic, error)
	GetByIDs(ids []uint) ([]*models.Analytic, error)
	Update(analytic *models.Analytic) error
	Delete(*models.Analytic) error
	DeleteByIDs(ids []uint) (err error)
	GetAll() ([]*models.Analytic, error)
	GetAllIds() ([]uint, error)
	WithPaginate(p contracts.Paginator, filter *models.AnalyticFilter) ([]*models.Analytic, error)
}

type repository struct {
	db *gorm.DB
	*app.Paginate
}

func NewAnalyticRepo(db *gorm.DB) AnalyticRepository {
	r := &repository{db: db}
	return r
}

func (r *repository) WithTx(tx *gorm.DB) AnalyticRepository {
	newR := &repository{db: tx}
	return newR
}

func (r *repository) Create(analytic *models.Analytic) error {
	return r.db.Create(analytic).Error
}

func (r *repository) GetByID(id uint) (*models.Analytic, error) {
	var analytic models.Analytic
	if err := r.db.First(&analytic, id).Error; err != nil {
		return nil, err
	}
	return &analytic, nil
}

func (r *repository) GetByIDs(ids []uint) ([]*models.Analytic, error) {
	var analytics []*models.Analytic
	if err := r.db.Where("id IN (?)", ids).Find(&analytics).Error; err != nil {
		return nil, err
	}
	return analytics, nil
}

func (r *repository) Update(analytic *models.Analytic) error {
	return r.db.Save(analytic).Error
}

func (r *repository) Delete(analytic *models.Analytic) error {
	return r.db.Delete(&models.Analytic{}, analytic.ID).Error
}

func (r *repository) DeleteByIDs(ids []uint) (err error) {
	return r.db.Where("id IN ?", ids).Delete(&models.Analytic{}).Error
}

func (r *repository) GetAll() ([]*models.Analytic, error) {
	var analytics []*models.Analytic
	if err := r.db.Order("id ASC").Find(&analytics).Error; err != nil {
		return nil, err
	}
	return analytics, nil
}

func (r *repository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.Analytic{}).Pluck("id", &ids).Error; err != nil {
		log.Printf("Failed to fetch IDs from the database: %v\n", err)
	}
	return ids, nil
}

func (r *repository) WithPaginate(paginator contracts.Paginator, filter *models.AnalyticFilter) ([]*models.Analytic, error) {
	var items []*models.Analytic
	var total int64

	model := models.Analytic{}
	table := model.GetTable()

	query := r.db.Model(&model)

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
		Order(fmt.Sprintf("%s.id ASC", model.GetTable())).
		Find(&items).Error
	if err != nil {
		return nil, err
	}

	paginator.SetTotal(int(total))
	return items, nil
}
