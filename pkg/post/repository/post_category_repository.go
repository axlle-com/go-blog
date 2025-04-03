package repository

import (
	"fmt"
	"github.com/axlle-com/blog/pkg/app/db"
	"github.com/axlle-com/blog/pkg/app/logger"
	app "github.com/axlle-com/blog/pkg/app/models"
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	"github.com/axlle-com/blog/pkg/post/models"
	"gorm.io/gorm"
	"strings"
)

type CategoryRepository interface {
	WithTx(tx *gorm.DB) CategoryRepository
	Create(postCategory *models.PostCategory) error
	GetByID(id uint) (*models.PostCategory, error)
	GetByIDs(ids []uint) ([]*models.PostCategory, error)
	Update(new *models.PostCategory, old *models.PostCategory) error
	DeleteByID(id uint) error
	Delete(category *models.PostCategory) error
	GetAll() ([]*models.PostCategory, error)
	GetAllForParent(parent *models.PostCategory) ([]*models.PostCategory, error)
	GetAllIds() ([]uint, error)
	WithPaginate(paginator contracts.Paginator, filter *models.CategoryFilter) ([]*models.PostCategory, error)
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

func (r *categoryRepository) GetDescendantsByID(id uint) ([]*models.PostCategory, error) {
	category, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	return r.GetDescendants(category)
}

func (r *categoryRepository) GetByIDs(ids []uint) ([]*models.PostCategory, error) {
	var categories []*models.PostCategory
	if err := r.db.Where("id IN (?)", ids).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) save(category *models.PostCategory) error {
	return r.db.Select(
		"UserID",
		"TemplateID",
		"PostCategoryID",
		"LeftSet",
		"RightSet",
		"Level",
		"MetaTitle",
		"MetaDescription",
		"Alias",
		"URL",
		"IsPublished",
		"IsFavourites",
		"InSitemap",
		"Image",
		"ShowImage",
		"Title",
		"TitleShort",
		"DescriptionPreview",
		"Description",
		"Sort",
	).Save(category).Error
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

func (r *categoryRepository) WithPaginate(p contracts.Paginator, filter *models.CategoryFilter) ([]*models.PostCategory, error) {
	var categories []*models.PostCategory
	var total int64

	category := models.PostCategory{}

	query := r.db.Model(&category)
	query.Count(&total)

	err := query.Scopes(r.SetPaginate(p.GetPage(), p.GetPageSize())).
		Order(fmt.Sprintf("%s.id ASC", category.GetTable())).
		Find(&categories).Error
	if err != nil {
		return nil, err
	}

	p.SetTotal(int(total))
	return categories, nil
}

func (r *categoryRepository) Delete(category *models.PostCategory) error {
	likePattern := fmt.Sprintf("%s%%", category.Path)
	return r.db.Where("path LIKE ?", likePattern).Delete(&models.PostCategory{}).Error
}

func (r *categoryRepository) GetRoots() ([]*models.PostCategory, error) {
	var roots []*models.PostCategory
	err := r.db.Where("post_category_id IS NULL").Find(&roots).Error
	return roots, err
}

func (r *categoryRepository) Create(category *models.PostCategory) error {
	category.Creating()
	if category.PostCategoryID == nil || *category.PostCategoryID == 0 {
		if err := r.db.Create(category).Error; err != nil {
			return err
		}
		// Для корневой категории путь – просто /id/
		category.Path = fmt.Sprintf("/%d/", category.ID)
		return r.db.Model(category).Update("path", category.Path).Error
	}

	// Если есть родитель, получаем его данные.
	var parent models.PostCategory
	if err := r.db.First(&parent, *category.PostCategoryID).Error; err != nil {
		return fmt.Errorf("не найден родитель: %w", err)
	}

	if err := r.db.Create(category).Error; err != nil {
		return err
	}
	// Путь дочернего узла – путь родителя + id дочернего.
	category.Path = fmt.Sprintf("%s%d/", parent.Path, category.ID)
	return r.db.Model(category).Update("path", category.Path).Error
}

func (r *categoryRepository) GetByID(id uint) (*models.PostCategory, error) {
	var category models.PostCategory
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetDescendants(category *models.PostCategory) ([]*models.PostCategory, error) {
	var descendants []*models.PostCategory
	likePattern := fmt.Sprintf("%s%%", category.Path)
	err := r.db.
		Where("path LIKE ? AND id <> ?", likePattern, category.ID).
		Order("id ASC").
		Find(&descendants).Error
	if err != nil {
		return nil, err
	}
	return descendants, nil
}

func (r *categoryRepository) GetAllForParent(parent *models.PostCategory) ([]*models.PostCategory, error) {
	var descendants []*models.PostCategory
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

func (r *categoryRepository) Update(new *models.PostCategory, old *models.PostCategory) error {
	new.Updating()

	// Если родитель не изменился – просто сохраняем изменения.
	oldParent, newParent := uint(0), uint(0)
	if old.PostCategoryID != nil {
		oldParent = *old.PostCategoryID
	}
	if new.PostCategoryID != nil {
		newParent = *new.PostCategoryID
	}

	if oldParent == newParent {
		return r.save(new)
	}

	// Если родитель меняется, требуется пересчитать путь для нового поддерева.
	var newParentCategory models.PostCategory
	if newParent != 0 {
		if err := r.db.First(&newParentCategory, newParent).Error; err != nil {
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
		newPath = fmt.Sprintf("%s%d/", newParentCategory.Path, new.ID)
	}

	// Обновляем путь для узла и всех потомков.
	var descendants []*models.PostCategory
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
