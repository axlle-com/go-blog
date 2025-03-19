package repository

import (
	"errors"
	"fmt"
	"github.com/axlle-com/blog/pkg/app/db"
	"github.com/axlle-com/blog/pkg/app/logger"
	app "github.com/axlle-com/blog/pkg/app/models"
	"github.com/axlle-com/blog/pkg/post/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	WithTx(tx *gorm.DB) CategoryRepository
	Create(postCategory *models.PostCategory) error
	GetByID(id uint) (*models.PostCategory, error)
	Update(new *models.PostCategory, old *models.PostCategory) error
	DeleteByID(id uint) error
	Delete(category *models.PostCategory) error
	GetAll() ([]*models.PostCategory, error)
	GetAllIds() ([]uint, error)
	WithPaginate(page, pageSize int) ([]*models.PostCategory, error)
	GetRoots() ([]*models.PostCategory, error)
	GetDescendants(category *models.PostCategory) ([]*models.PostCategory, error)
	GetDescendantsByID(id uint) ([]*models.PostCategory, error)
}

type categoryRepository struct {
	db *gorm.DB
	*app.Paginate
}

func NewCategoryRepo() CategoryRepository {
	return &categoryRepository{db: db.GetDB()}
}

func (r *categoryRepository) WithTx(tx *gorm.DB) CategoryRepository {
	return &categoryRepository{db: tx}
}

func (r *categoryRepository) GetDescendants(category *models.PostCategory) ([]*models.PostCategory, error) {
	var descendants []*models.PostCategory
	err := r.db.
		Where("left_set > ? AND right_set < ?", category.LeftSet, category.RightSet).
		Order("left_set ASC").
		Find(&descendants).Error
	if err != nil {
		return nil, err
	}
	return descendants, nil
}

func (r *categoryRepository) GetDescendantsByID(id uint) ([]*models.PostCategory, error) {
	category, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	return r.GetDescendants(category)
}

func (r *categoryRepository) Create(category *models.PostCategory) error {
	category.Creating()

	if category.PostCategoryID == nil || *category.PostCategoryID == 0 {
		return r.createRoot(category)
	}
	return r.createChild(category)
}

func (r *categoryRepository) Delete(category *models.PostCategory) error {
	var node models.PostCategory
	if err := r.db.First(&node, category.ID).Error; err != nil {
		return err
	}

	size := node.RightSet - node.LeftSet + 1

	// Удаляем узел и потомков
	if err := r.db.Where("left_set BETWEEN ? AND ?", node.LeftSet, node.RightSet).
		Delete(&models.PostCategory{}).Error; err != nil {
		return err
	}

	// Корректируем только если это не последний корень
	if !(node.PostCategoryID == nil || *node.PostCategoryID == 0) {
		if err := r.db.Model(&models.PostCategory{}).
			Where("right_set > ?", node.RightSet).
			UpdateColumn("right_set", gorm.Expr("right_set - ?", size)).
			Error; err != nil {
			return err
		}

		return r.db.Model(&models.PostCategory{}).
			Where("left_set > ?", node.RightSet).
			UpdateColumn("left_set", gorm.Expr("left_set - ?", size)).
			Error
	}

	return nil
}

func (r *categoryRepository) GetRoots() ([]*models.PostCategory, error) {
	var roots []*models.PostCategory
	err := r.db.Where("post_category_id IS NULL").Find(&roots).Error
	return roots, err
}

func (r *categoryRepository) GetByID(id uint) (*models.PostCategory, error) {
	var postCategory models.PostCategory
	if err := r.db.First(&postCategory, id).Error; err != nil {
		return nil, err
	}
	return &postCategory, nil
}

func (r *categoryRepository) Update(new *models.PostCategory, old *models.PostCategory) error {
	new.Updating()
	oldID, newID := uint(0), uint(0)
	if old.PostCategoryID != nil {
		oldID = *old.PostCategoryID
	}
	if new.PostCategoryID != nil {
		newID = *new.PostCategoryID
	}
	if oldID != newID {
		new.LeftSet = old.LeftSet
		new.RightSet = old.RightSet
		new.Level = old.Level

		err := r.moveTo(new, new.PostCategoryID)
		if err != nil {
			return err
		}
		return nil
	}

	return r.db.Save(new).Error
}

func (r *categoryRepository) DeleteByID(id uint) error {
	var node models.PostCategory
	if err := r.db.First(&node, id).Error; err != nil {
		return err
	}
	return r.Delete(&node)
}

func (r *categoryRepository) GetAll() ([]*models.PostCategory, error) {
	var postCategories []*models.PostCategory
	if err := r.db.Find(&postCategories).Error; err != nil {
		return nil, err
	}
	return postCategories, nil
}

func (r *categoryRepository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.PostCategory{}).Pluck("id", &ids).Error; err != nil {
		logger.Error(err)
	}
	return ids, nil
}

func (r *categoryRepository) WithPaginate(page, pageSize int) ([]*models.PostCategory, error) {
	var categories []*models.PostCategory

	err := r.db.Model(&models.PostCategory{}).Scopes(r.SetPaginate(page, pageSize)).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) createRoot(category *models.PostCategory) error {
	var maxRight struct{ Value int }
	if err := r.db.Model(&models.PostCategory{}).
		Select("COALESCE(MAX(right_set), 0) as value").
		Scan(&maxRight).Error; err != nil {
		return err
	}

	category.LeftSet = maxRight.Value + 1
	category.RightSet = maxRight.Value + 2
	category.Level = 0
	category.PostCategoryID = nil

	return r.db.Create(category).Error
}

func (r *categoryRepository) createChild(category *models.PostCategory) error {
	var parent models.PostCategory
	if category.PostCategoryID == nil {
		return errors.New("parent category required for child")
	}

	if err := r.db.First(&parent, *category.PostCategoryID).Error; err != nil {
		return err
	}

	// Сохраняем исходное значение правой границы родителя
	target := parent.RightSet

	// Обновляем правую границу для всех узлов, у которых right_set >= target
	if err := r.db.Model(&models.PostCategory{}).
		Where("right_set >= ?", target).
		UpdateColumn("right_set", gorm.Expr("right_set + 2")).
		Error; err != nil {
		return err
	}

	// Обновляем левую границу для всех узлов, у которых left_set > target
	if err := r.db.Model(&models.PostCategory{}).
		Where("left_set > ?", target).
		UpdateColumn("left_set", gorm.Expr("left_set + 2")).
		Error; err != nil {
		return err
	}

	// Назначаем границы для нового узла, используя сохранённое значение target
	category.LeftSet = target
	category.RightSet = target + 1
	category.Level = parent.Level + 1

	return r.db.Create(category).Error
}

func (r *categoryRepository) moveTo(category *models.PostCategory, newParentID *uint) error {
	// Перемещение в корень
	if newParentID == nil {
		return r.moveToRoot(category)
	}

	// Стандартное перемещение
	return r.moveToParent(category, *newParentID)
}

func (r *categoryRepository) moveToRoot(category *models.PostCategory) error {
	var maxRight struct{ Value int }
	if err := r.db.Model(&models.PostCategory{}).
		Select("COALESCE(MAX(right_set), 0) as value").
		Scan(&maxRight).Error; err != nil {
		return fmt.Errorf("failed to get max right_set: %w", err)
	}

	// Рассчитываем параметры для перемещения
	size := category.RightSet - category.LeftSet + 1
	newLeft := maxRight.Value + 1
	newRight := maxRight.Value + size
	offset := newLeft - category.LeftSet

	// Обновляем границы перемещаемого поддерева
	if err := r.db.Model(&models.PostCategory{}).
		Where("left_set BETWEEN ? AND ?", category.LeftSet, category.RightSet).
		Updates(map[string]interface{}{
			"left_set":  gorm.Expr("left_set + ?", offset),
			"right_set": gorm.Expr("right_set + ?", offset),
			"level":     gorm.Expr("level - ?", category.Level),
		}).Error; err != nil {
		return fmt.Errorf("failed to update subtree: %w", err)
	}

	// Корректируем оставшиеся узлы
	if err := r.db.Model(&models.PostCategory{}).
		Where("left_set > ?", category.RightSet).
		UpdateColumn("left_set", gorm.Expr("left_set - ?", size)).
		Error; err != nil {
		return fmt.Errorf("failed to shift left_sets: %w", err)
	}

	if err := r.db.Model(&models.PostCategory{}).
		Where("right_set > ?", category.RightSet).
		UpdateColumn("right_set", gorm.Expr("right_set - ?", size)).
		Error; err != nil {
		return fmt.Errorf("failed to shift right_sets: %w", err)
	}

	// Обновляем параметры корневой категории
	category.PostCategoryID = nil
	category.LeftSet = newLeft
	category.RightSet = newRight
	category.Level = 0

	return r.db.Save(category).Error
}

func (r *categoryRepository) moveToParent(category *models.PostCategory, newParentID uint) error {
	var newParent models.PostCategory
	if err := r.db.First(&newParent, newParentID).Error; err != nil {
		return fmt.Errorf("new parent not found: %w", err)
	}

	// Нельзя перемещать узел в самого себя или в своего потомка
	if category.ID == newParent.ID {
		return errors.New("cannot move to self")
	}
	if newParent.LeftSet > category.LeftSet && newParent.RightSet < category.RightSet {
		return errors.New("cannot move to descendant")
	}

	size := category.RightSet - category.LeftSet + 1
	newPosition := newParent.RightSet

	// Освобождаем место в новом положении
	if err := r.db.Model(&models.PostCategory{}).
		Where("right_set >= ?", newPosition).
		UpdateColumn("right_set", gorm.Expr("right_set + ?", size)).
		Error; err != nil {
		return err
	}

	if err := r.db.Model(&models.PostCategory{}).
		Where("left_set >= ?", newPosition).
		UpdateColumn("left_set", gorm.Expr("left_set + ?", size)).
		Error; err != nil {
		return err
	}

	// Если узел перемещается вправо от исходной позиции
	if category.LeftSet < newPosition {
		newPosition -= size
	}

	// Вычисляем смещения для обновления
	offset := newPosition - category.LeftSet
	levelDiff := newParent.Level + 1 - category.Level

	// Обновляем перемещаемый узел и его потомков в базе
	if err := r.db.Model(&models.PostCategory{}).
		Where("left_set BETWEEN ? AND ?", category.LeftSet, category.RightSet).
		Updates(map[string]interface{}{
			"left_set":  gorm.Expr("left_set + ?", offset),
			"right_set": gorm.Expr("right_set + ?", offset),
			"level":     gorm.Expr("level + ?", levelDiff),
		}).Error; err != nil {
		return err
	}

	// Обновляем значения в структуре, чтобы они соответствовали изменениям в базе
	category.LeftSet += offset
	category.RightSet += offset
	category.Level += levelDiff

	// Заполняем освободившееся место в базе
	if err := r.db.Model(&models.PostCategory{}).
		Where("left_set > ?", category.RightSet).
		UpdateColumn("left_set", gorm.Expr("left_set - ?", size)).
		Error; err != nil {
		return err
	}

	if err := r.db.Model(&models.PostCategory{}).
		Where("right_set > ?", category.RightSet).
		UpdateColumn("right_set", gorm.Expr("right_set - ?", size)).
		Error; err != nil {
		return err
	}

	// Обновляем ссылку на родителя в структуре
	category.PostCategoryID = &newParent.ID
	return r.db.Save(category).Error
}
