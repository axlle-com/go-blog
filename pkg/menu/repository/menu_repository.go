package repository

import (
	"errors"
	"fmt"

	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/menu/models"
	"gorm.io/gorm"
)

type MenuRepository interface {
	WithTx(tx *gorm.DB) MenuRepository
	Create(menu *models.Menu) error
	GetByID(id uint) (*models.Menu, error)
	GetByParam(field string, value any) (*models.Menu, error)
	GetByParams(params map[string]any) ([]*models.Menu, error)
	Update(menu *models.Menu) error
	Delete(menu *models.Menu) error
	GetAll() ([]*models.Menu, error)
	GetAllIds() ([]uint, error)
	WithPaginate(paginator contract.Paginator, filter *models.MenuFilter) ([]*models.Menu, error)
}

type menuRepository struct {
	db *gorm.DB
	*app.Paginate
}

func NewMenuRepo(db *gorm.DB) MenuRepository {
	r := &menuRepository{db: db}
	return r
}

func (r *menuRepository) WithTx(tx *gorm.DB) MenuRepository {
	return &menuRepository{db: tx}
}

func (r *menuRepository) Create(menu *models.Menu) error {
	menu.Creating()
	return r.db.Create(menu).Error
}

func (r *menuRepository) GetByID(id uint) (*models.Menu, error) {
	var model models.Menu
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *menuRepository) Update(menu *models.Menu) error {
	menu.Updating()
	return r.db.Select(
		"TemplateID",
		"Title",
		"IsPublished",
		"IsMain",
		"Ico",
		"Sort",
	).Save(menu).Error
}

func (r *menuRepository) Delete(menu *models.Menu) error {
	if menu.Deleting() {
		return r.db.Delete(&models.Menu{}, menu.ID).Error
	}
	return errors.New("deletion errors occurred")
}

func (r *menuRepository) GetAll() ([]*models.Menu, error) {
	var menus []*models.Menu
	if err := r.db.Order("id ASC").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *menuRepository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.Menu{}).Pluck("id", &ids).Error; err != nil {
		logger.Error(err)
	}
	return ids, nil
}

func (r *menuRepository) WithPaginate(p contract.Paginator, filter *models.MenuFilter) ([]*models.Menu, error) {
	var menus []*models.Menu
	var total int64

	menu := models.Menu{}
	table := menu.GetTable()

	query := r.db.Model(&menus)

	// TODO WHERE IN; LIKE
	for col, val := range filter.GetMap() {
		if col == "title" {
			query = query.Where(fmt.Sprintf("%s.%v ilike ?", table, col), fmt.Sprintf("%%%v%%", val))
			continue
		}
		query = query.Where(fmt.Sprintf("%s.%v = ?", table, col), val)
	}

	query.Count(&total)

	err := query.Scopes(r.SetPaginate(p.GetPage(), p.GetPageSize())).
		Order(fmt.Sprintf("%s.id ASC", table)).
		Scan(&menus).Error

	p.SetTotal(int(total))
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *menuRepository) GetByParam(field string, value any) (*models.Menu, error) {
	var menu models.Menu
	condition := map[string]any{
		field: value,
	}
	if err := r.db.Where(condition).First(&menu).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *menuRepository) GetByParams(params map[string]any) ([]*models.Menu, error) {
	var menus []*models.Menu
	if err := r.db.Where(params).Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}
