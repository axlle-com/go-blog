package repository

import (
	"fmt"
	"strings"

	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MenuItemRepository interface {
	WithTx(tx *gorm.DB) MenuItemRepository
	Create(menuItem *models.MenuItem) error
	GetByID(id uint) (*models.MenuItem, error)
	GetByIDs(ids []uint) ([]*models.MenuItem, error)
	Update(new *models.MenuItem, old *models.MenuItem) error
	UpdateURLForPublisher(publisherUuid uuid.UUID, newURL string) (int64, error)
	DetachPublisher(publisherUuid uuid.UUID) (int64, error)
	DeleteByID(id uint) error
	Delete(menuItem *models.MenuItem) error
	GetByFilter(p contract.Paginator, filter *models.MenuItemFilter) ([]*models.MenuItem, error)
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
	var items []*models.MenuItem
	if err := r.db.Where("id IN (?)", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *menuItemRepository) save(menuItem *models.MenuItem) error {
	return r.db.Select(menuItem.Fields()).Save(menuItem).Error
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
	if err := r.db.
		Where(params).
		Order(fmt.Sprintf("%s.sort ASC", (&models.MenuItem{}).GetTable())).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *menuItemRepository) GetByFilter(paginator contract.Paginator, filter *models.MenuItemFilter) ([]*models.MenuItem, error) {
	var items []*models.MenuItem
	model := models.MenuItem{}
	var total int64

	query := r.db.Model(&model)

	if filter != nil {
		if filter.MenuID != nil && *filter.MenuID != 0 {
			query = query.Where("menu_id = ?", *filter.MenuID)
		}

		// исключить узел и всех его потомков из выборки (чтобы нельзя было выбрать родителем потомка)
		if filter.ForNotMenuItemID != nil && *filter.ForNotMenuItemID != 0 {
			var nodePathLtree string
			err := r.db.Table(model.GetTable()).
				Where("id = ?", *filter.ForNotMenuItemID).
				Pluck("path_ltree", &nodePathLtree).Error
			if err != nil {
				return nil, err
			}
			if nodePathLtree != "" {
				query = query.Where("NOT (path_ltree <@ ?::ltree)", nodePathLtree)
			}
		}

		if filter.IDs != nil {
			query = query.Where("id IN ?", filter.IDs)
		}
	}

	if paginator == nil {
		err := query.Order(fmt.Sprintf("%s.sort ASC", model.GetTable())).Find(&items).Error

		return items, err
	}

	query.Count(&total)

	err := query.Scopes(r.SetPaginate(paginator.GetPage(), paginator.GetPageSize())).
		Order(fmt.Sprintf("%s.sort ASC", model.GetTable())).
		Find(&items).Error
	if err != nil {
		return nil, err
	}

	paginator.SetTotal(int(total))
	return items, nil
}

func (r *menuItemRepository) GetAll() ([]*models.MenuItem, error) {
	var items []*models.MenuItem
	if err := r.db.Order("id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *menuItemRepository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.MenuItem{}).Pluck("id", &ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *menuItemRepository) Delete(menuItem *models.MenuItem) error {
	if menuItem == nil {
		return fmt.Errorf("menuItem is nil")
	}

	// если PathLtree не заполнен — подстрахуемся и возьмём из БД
	path := menuItem.PathLtree
	if path == "" && menuItem.ID != 0 {
		var tmp models.MenuItem
		if err := r.db.Select("path_ltree").First(&tmp, menuItem.ID).Error; err != nil {
			return err
		}

		path = tmp.PathLtree
	}

	if path == "" {
		return fmt.Errorf("empty path_ltree for delete")
	}

	return r.db.Where("path_ltree <@ ?::ltree", path).Delete(&models.MenuItem{}).Error
}

func (r *menuItemRepository) GetRoots() ([]*models.MenuItem, error) {
	var roots []*models.MenuItem
	err := r.db.Where("menu_item_id IS NULL").Find(&roots).Error

	return roots, err
}

func (r *menuItemRepository) Create(menuItem *models.MenuItem) error {
	if menuItem == nil {
		return fmt.Errorf("menuItem is nil")
	}

	menuItem.Creating()

	// Если корень (MenuItemID NULL или 0)
	if menuItem.MenuItemID == nil || *menuItem.MenuItemID == 0 {
		if err := r.db.Create(menuItem).Error; err != nil {
			return err
		}

		menuItem.PathLtree = fmt.Sprintf("%d", menuItem.ID)
		return r.db.Model(menuItem).Update("path_ltree", menuItem.PathLtree).Error
	}

	// родитель
	var parent models.MenuItem
	if err := r.db.First(&parent, *menuItem.MenuItemID).Error; err != nil {
		return fmt.Errorf("parent not found: %w", err)
	}
	if parent.PathLtree == "" {
		return fmt.Errorf("parent has empty path_ltree (id=%d)", parent.ID)
	}

	if err := r.db.Create(menuItem).Error; err != nil {
		return err
	}

	menuItem.PathLtree = fmt.Sprintf("%s.%d", parent.PathLtree, menuItem.ID)

	return r.db.Model(menuItem).Update("path_ltree", menuItem.PathLtree).Error
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
	err := r.db.
		Where("path_ltree <@ ?::ltree AND id <> ?", menuItem.PathLtree, menuItem.ID).
		Order("id ASC").
		Find(&descendants).Error
	if err != nil {
		return nil, err
	}
	return descendants, nil
}

func (r *menuItemRepository) GetAllForParent(parent *models.MenuItem) ([]*models.MenuItem, error) {
	var items []*models.MenuItem
	err := r.db.
		Where("NOT (path_ltree <@ ?::ltree) AND id <> ?", parent.PathLtree, parent.ID).
		Order("id ASC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func validatePathEndsWithID(path string, id uint) error {
	if path == "" {
		return fmt.Errorf("empty path_ltree")
	}

	last := path
	if i := strings.LastIndex(path, "."); i >= 0 {
		last = path[i+1:]
	}
	want := fmt.Sprintf("%d", id)
	if last != want {
		return fmt.Errorf("bad path_ltree=%q: last label=%q, expected id=%q", path, last, want)
	}

	return nil
}

func (r *menuItemRepository) Update(new *models.MenuItem, _ *models.MenuItem) error {
	if new == nil {
		return fmt.Errorf("new is nil")
	}

	new.Updating()

	return r.db.Transaction(func(db *gorm.DB) error {
		// 1) читаем old под блокировкой
		var old models.MenuItem
		if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&old, new.ID).Error; err != nil {
			return err
		}

		if err := validatePathEndsWithID(old.PathLtree, old.ID); err != nil {
			return fmt.Errorf("refuse to update because old path_ltree is invalid: %w", err)
		}

		// 2) определяем oldParent/newParent
		oldParent, newParent := uint(0), uint(0)
		if old.MenuItemID != nil {
			oldParent = *old.MenuItemID
		}
		if new.MenuItemID != nil {
			newParent = *new.MenuItemID
		}

		// 2.0) защита: parent=self
		if newParent != 0 && newParent == old.ID {
			return fmt.Errorf("cannot set parent to self: id=%d", old.ID)
		}

		// ВАЖНО: не даём случайно менять menu_id через Update (иначе можно разорвать целостность дерева)
		new.MenuID = old.MenuID

		// 2.1) если родитель не менялся — обычный save, но path_ltree оставляем прежним
		if oldParent == newParent {
			new.PathLtree = old.PathLtree
			return db.Select(new.Fields()).Save(new).Error
		}

		// 3) валидируем нового родителя (если перенос не в корень)
		var newParentMenuItem models.MenuItem
		if newParent != 0 {
			if err := db.Select("id", "path_ltree", "menu_id").
				First(&newParentMenuItem, newParent).Error; err != nil {
				return fmt.Errorf("new parent not found: %w", err)
			}
			if newParentMenuItem.PathLtree == "" {
				return fmt.Errorf("new parent has empty path_ltree (id=%d)", newParentMenuItem.ID)
			}
			// защита: нельзя переносить в другой menu
			if newParentMenuItem.MenuID != old.MenuID {
				return fmt.Errorf(
					"cannot move node between menus: old menu_id=%d new parent menu_id=%d",
					old.MenuID, newParentMenuItem.MenuID,
				)
			}

			// защита: нельзя завернуть в своего потомка (строковая проверка по ltree-путям)
			// (newParentMenuItem.PathLtree начинается с old.PathLtree+".")
			oldPrefix := old.PathLtree
			descPrefix := oldPrefix + "."
			if newParentMenuItem.PathLtree == oldPrefix || strings.HasPrefix(newParentMenuItem.PathLtree, descPrefix) {
				return fmt.Errorf(
					"cannot move node under its descendant: node=%d new_parent=%d",
					old.ID, newParentMenuItem.ID,
				)
			}
		}

		oldPathLtree := old.PathLtree

		var newPathLtree string
		if newParent == 0 {
			newPathLtree = fmt.Sprintf("%d", new.ID)
		} else {
			newPathLtree = fmt.Sprintf("%s.%d", newParentMenuItem.PathLtree, new.ID)
		}

		// 4) обновляем потомков (кроме корня поддерева), строго внутри menu_id
		if err := db.Exec(`
			UPDATE menu_items
			SET path_ltree = ?::ltree || subpath(path_ltree, nlevel(?::ltree))
			WHERE menu_id = ?
			  AND path_ltree <@ ?::ltree
			  AND path_ltree <> ?::ltree
		`, newPathLtree, oldPathLtree, old.MenuID, oldPathLtree, oldPathLtree).Error; err != nil {
			return fmt.Errorf("failed to update subtree paths: %w", err)
		}

		// 5) обновляем сам узел
		new.PathLtree = newPathLtree
		return db.Select(new.Fields()).Save(new).Error
	})
}

func (r *menuItemRepository) UpdateURLForPublisher(publisherUuid uuid.UUID, newURL string) (int64, error) {
	res := r.db.Model(&models.MenuItem{}).
		Where("publisher_uuid = ?", publisherUuid).
		Updates(map[string]any{"url": newURL})

	return res.RowsAffected, res.Error
}

func (r *menuItemRepository) DetachPublisher(publisherUuid uuid.UUID) (int64, error) {
	res := r.db.Model(&models.MenuItem{}).
		Where("publisher_uuid = ?", publisherUuid).
		Updates(map[string]any{
			"publisher_uuid": gorm.Expr("NULL"),
			"url":            "#",
		})
	return res.RowsAffected, res.Error
}
