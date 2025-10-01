package repository

import (
	"fmt"
	"strings"

	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/menu/models"
	"gorm.io/gorm"
)

type MenuItemRepository interface {
	WithTx(tx *gorm.DB) MenuItemRepository
	Create(menuItem *models.MenuItem) error
	GetByID(id uint) (*models.MenuItem, error)
	GetByIDs(ids []uint) ([]*models.MenuItem, error)
	Update(new *models.MenuItem, old *models.MenuItem) error
	DeleteByID(id uint) error
	Delete(menuItem *models.MenuItem) error
	GetByFilter(p contracts.Paginator, filter *models.MenuItemFilter) ([]*models.MenuItem, error)
	GetByParams(params map[string]any) ([]*models.MenuItem, error)
	GetAll() ([]*models.MenuItem, error)
	GetAllForParent(parent *models.MenuItem) ([]*models.MenuItem, error)
	GetAllIds() ([]uint, error)
	GetRoots() ([]*models.MenuItem, error)
	GetDescendants(menuItem *models.MenuItem) ([]*models.MenuItem, error)
	GetDescendantsByID(id uint) ([]*models.MenuItem, error)
}

type menuItemRepository struct {
	db *gorm.DB
	*app.Paginate
}

func NewMenuItemRepo(db *gorm.DB) MenuItemRepository {
	return &menuItemRepository{db: db}
}

func (r *menuItemRepository) WithTx(tx *gorm.DB) MenuItemRepository {
	return &menuItemRepository{db: tx}
}

func (r *menuItemRepository) GetDescendantsByID(id uint) ([]*models.MenuItem, error) {
	menuItem, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	return r.GetDescendants(menuItem)
}

func (r *menuItemRepository) GetByIDs(ids []uint) ([]*models.MenuItem, error) {
	var categories []*models.MenuItem
	if err := r.db.Where("id IN (?)", ids).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *menuItemRepository) save(menuItem *models.MenuItem) error {
	return r.db.Select(
		"UserID",
		"PublisherUUID",
		"MenuID",
		"MenuItemID",
		"Path",
		"Title",
		"URL",
		"IsPublished",
		"Ico",
		"Sort",
	).Save(menuItem).Error
}

func (r *menuItemRepository) DeleteByID(id uint) error {
	var node models.MenuItem
	if err := r.db.First(&node, id).Error; err != nil {
		return err
	}
	return r.Delete(&node)
}

func (r *menuItemRepository) GetByParams(params map[string]any) ([]*models.MenuItem, error) {
	var items []*models.MenuItem
	if err := r.db.Where(params).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *menuItemRepository) GetByFilter(paginator contracts.Paginator, filter *models.MenuItemFilter) ([]*models.MenuItem, error) {
	var items []*models.MenuItem
	model := models.MenuItem{}
	var total int64

	query := r.db.Model(&model)

	if filter.MenuID != nil && *filter.MenuID != 0 {
		query = query.Where("menu_id = ?", *filter.MenuID)
	}

	if filter.ForNotMenuItemID != nil && *filter.ForNotMenuItemID != 0 {
		var nodePath string
		err := r.db.Table(model.GetTable()).
			Where("id = ?", *filter.ForNotMenuItemID).
			Pluck("path", &nodePath).Error
		if err != nil {
			return nil, err
		}

		query = query.Where("path NOT LIKE ?", nodePath+"%")
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

func (r *menuItemRepository) GetAll() ([]*models.MenuItem, error) {
	var menuItems []*models.MenuItem
	if err := r.db.Order("id ASC").Find(&menuItems).Error; err != nil {
		return nil, err
	}
	return menuItems, nil
}

func (r *menuItemRepository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.MenuItem{}).Pluck("id", &ids).Error; err != nil {
		logger.Error(err)
	}
	return ids, nil
}

func (r *menuItemRepository) Delete(menuItem *models.MenuItem) error {
	likePattern := fmt.Sprintf("%s%%", menuItem.Path)
	return r.db.Where("path LIKE ?", likePattern).Delete(&models.MenuItem{}).Error
}

func (r *menuItemRepository) GetRoots() ([]*models.MenuItem, error) {
	var roots []*models.MenuItem
	err := r.db.Where("menu_item_id IS NULL").Find(&roots).Error
	return roots, err
}

func (r *menuItemRepository) Create(menuItem *models.MenuItem) error {
	menuItem.Creating()
	if menuItem.MenuItemID == nil || *menuItem.MenuItemID == 0 {
		if err := r.db.Create(menuItem).Error; err != nil {
			return err
		}
		// Для корневой категории путь – просто /id/
		menuItem.Path = fmt.Sprintf("/%d/", menuItem.ID)
		return r.db.Model(menuItem).Update("path", menuItem.Path).Error
	}

	// Если есть родитель, получаем его данные.
	var parent models.MenuItem
	if err := r.db.First(&parent, *menuItem.MenuItemID).Error; err != nil {
		return fmt.Errorf("не найден родитель: %w", err)
	}

	if err := r.db.Create(menuItem).Error; err != nil {
		return err
	}
	// Путь дочернего узла – путь родителя + id дочернего.
	menuItem.Path = fmt.Sprintf("%s%d/", parent.Path, menuItem.ID)
	return r.db.Model(menuItem).Update("path", menuItem.Path).Error
}

func (r *menuItemRepository) GetByID(id uint) (*models.MenuItem, error) {
	var menuItem models.MenuItem
	if err := r.db.First(&menuItem, id).Error; err != nil {
		return nil, err
	}
	return &menuItem, nil
}

func (r *menuItemRepository) GetDescendants(menuItem *models.MenuItem) ([]*models.MenuItem, error) {
	var descendants []*models.MenuItem
	likePattern := fmt.Sprintf("%s%%", menuItem.Path)
	err := r.db.
		Where("path LIKE ? AND id <> ?", likePattern, menuItem.ID).
		Order("id ASC").
		Find(&descendants).Error
	if err != nil {
		return nil, err
	}
	return descendants, nil
}

func (r *menuItemRepository) GetAllForParent(parent *models.MenuItem) ([]*models.MenuItem, error) {
	var descendants []*models.MenuItem
	likePattern := fmt.Sprintf("%s%%", parent.Path)
	err := r.db.
		Where("path NOT LIKE ? AND id <> ?", likePattern, parent.ID).
		Order("id ASC").
		Find(&descendants).Error
	if err != nil {
		return nil, err
	}
	return descendants, nil
}

func (r *menuItemRepository) Update(new *models.MenuItem, old *models.MenuItem) error {
	new.Updating()

	// Если родитель не изменился – просто сохраняем изменения.
	oldParent, newParent := uint(0), uint(0)
	if old.MenuItemID != nil {
		oldParent = *old.MenuItemID
	}
	if new.MenuItemID != nil {
		newParent = *new.MenuItemID
	}

	if oldParent == newParent {
		return r.save(new)
	}

	// Если родитель меняется, требуется пересчитать путь для нового поддерева.
	var newParentMenuItem models.MenuItem
	if newParent != 0 {
		if err := r.db.First(&newParentMenuItem, newParent).Error; err != nil {
			return fmt.Errorf("не найден новый родитель: %w", err)
		}
	}

	// Сохраняем старый путь для поиска потомков.
	oldPath := old.Path
	// Сохраняем новый путь для узла.
	var newPath string
	if newParent == 0 {
		// Перемещение в корень
		newPath = fmt.Sprintf("/%d/", new.ID)
	} else {
		newPath = fmt.Sprintf("%s%d/", newParentMenuItem.Path, new.ID)
	}

	// Обновляем путь для узла и всех потомков.
	var descendants []*models.MenuItem
	if err := r.db.
		Where("path LIKE ?", fmt.Sprintf("%s%%", oldPath)).
		Find(&descendants).Error; err != nil {
		return err
	}

	// Рассчитываем смещение нового пути относительно старого.
	// Для каждого потомка новый путь = newPath + (old descendant.Path без префикса oldPath).
	for _, node := range descendants {
		relative := strings.TrimPrefix(node.Path, oldPath)
		node.Path = newPath + relative
		if err := r.db.Model(node).Update("path", node.Path).Error; err != nil {
			return err
		}
	}

	// Обновляем сам узел.
	new.Path = newPath
	return r.save(new)
}
