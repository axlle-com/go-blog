package repository

import (
	"fmt"
	"strings"

	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CategoryRepository interface {
	Tx() *gorm.DB
	WithTx(tx *gorm.DB) CategoryRepository
	Create(postCategory *models.PostCategory) error
	GetByID(id uint) (*models.PostCategory, error)
	GetByIDs(ids []uint) ([]*models.PostCategory, error)
	FindByParam(field string, value any) (*models.PostCategory, error)
	Update(new *models.PostCategory, old *models.PostCategory) error
	UpdateFieldsByUUIDs(uuids []uuid.UUID, patch map[string]any) (int64, error)
	DeleteByID(id uint) error
	Delete(category *models.PostCategory) error
	GetAll() ([]*models.PostCategory, error)
	GetAllForParent(parent *models.PostCategory) ([]*models.PostCategory, error)
	GetAllIds() ([]uint, error)
	WithPaginate(paginator contract.Paginator, filter *models.CategoryFilter) ([]*models.PostCategory, error)
	GetRoots() ([]*models.PostCategory, error)
	GetDescendants(category *models.PostCategory) ([]*models.PostCategory, error)
	GetDescendantsByID(id uint) ([]*models.PostCategory, error)
}

type categoryRepository struct {
	db *gorm.DB
	*app.Paginate
}

func NewCategoryRepo(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) WithTx(tx *gorm.DB) CategoryRepository {
	return &categoryRepository{db: tx}
}

func (r *categoryRepository) Tx() *gorm.DB {
	return r.db.Begin()
}

func (r *categoryRepository) GetDescendantsByID(id uint) ([]*models.PostCategory, error) {
	category, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	return r.GetDescendants(category)
}

func (r *categoryRepository) GetByIDs(ids []uint) ([]*models.PostCategory, error) {
	if len(ids) == 0 {
		return []*models.PostCategory{}, nil
	}

	var categories []*models.PostCategory
	if err := r.db.Where("id IN ?", ids).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) save(category *models.PostCategory) error {
	return r.db.Select(category.Fields()).Save(category).Error
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
	if err := r.db.Order("id ASC").Find(&postCategories).Error; err != nil {
		return nil, err
	}

	return postCategories, nil
}

func (r *categoryRepository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.PostCategory{}).Pluck("id", &ids).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	return ids, nil
}

func (r *categoryRepository) WithPaginate(p contract.Paginator, filter *models.CategoryFilter) ([]*models.PostCategory, error) {
	var categories []*models.PostCategory
	var total int64

	category := models.PostCategory{}
	table := category.GetTable()

	query := r.db.Model(&category)

	if filter != nil {
		if filter.ID != nil {
			query = query.Where(fmt.Sprintf("%s.id = ?", table), *filter.ID)
		}
		if filter.TemplateName != nil && *filter.TemplateName != "" {
			query = query.Where(fmt.Sprintf("%s.template_name = ?", table), *filter.TemplateName)
		}
		if filter.UserID != nil {
			query = query.Where(fmt.Sprintf("%s.user_id = ?", table), *filter.UserID)
		}
		if filter.PostCategoryID != nil {
			query = query.Where(fmt.Sprintf("%s.post_category_id = ?", table), *filter.PostCategoryID)
		}
		if filter.Title != nil && *filter.Title != "" {
			query = query.Where(fmt.Sprintf("%s.title ilike ?", table), fmt.Sprintf("%%%s%%", *filter.Title))
		}
		if filter.Query != nil && *filter.Query != "" {
			query = query.Where(fmt.Sprintf("%s.title ilike ?", table), fmt.Sprintf("%%%s%%", *filter.Query))
		}
		if filter.URL != nil && *filter.URL != "" {
			query = query.Where(fmt.Sprintf("%s.url = ?", table), *filter.URL)
		}
		if filter.Date != nil && *filter.Date != "" {
			query = query.Where(fmt.Sprintf("DATE(%s.created_at) = ?", table), *filter.Date)
		}
		if len(filter.UUIDs) > 0 {
			query = query.Where(fmt.Sprintf("%s.uuid IN ?", table), filter.UUIDs)
		}
	}

	query.Count(&total)

	err := query.Scopes(r.SetPaginate(p.GetPage(), p.GetPageSize())).
		Order(fmt.Sprintf("%s.id ASC", table)).
		Find(&categories).Error
	if err != nil {
		return nil, err
	}

	p.SetTotal(int(total))
	return categories, nil
}

func (r *categoryRepository) Delete(category *models.PostCategory) error {
	if category == nil {
		return fmt.Errorf("category is nil")
	}

	// если PathLtree не заполнен — подстрахуемся и возьмём из БД
	path := category.PathLtree
	if path == "" && category.ID != 0 {
		var tmp models.PostCategory
		if err := r.db.Select("path_ltree").First(&tmp, category.ID).Error; err != nil {
			return err
		}

		path = tmp.PathLtree
	}

	if path == "" {
		return fmt.Errorf("empty path_ltree for delete (id=%d)", category.ID)
	}

	// удаляем узел и всех потомков
	return r.db.Where("path_ltree <@ ?::ltree", path).
		Delete(&models.PostCategory{}).Error
}

func (r *categoryRepository) GetRoots() ([]*models.PostCategory, error) {
	var roots []*models.PostCategory
	err := r.db.Where("post_category_id IS NULL").Find(&roots).Error
	return roots, err
}

func (r *categoryRepository) Create(category *models.PostCategory) error {
	if category == nil {
		return fmt.Errorf("category is nil")
	}

	category.Creating()

	// root
	if category.PostCategoryID == nil || *category.PostCategoryID == 0 {
		if err := r.db.Create(category).Error; err != nil {
			return err
		}
		category.PathLtree = fmt.Sprintf("%d", category.ID)
		return r.db.Model(category).UpdateColumn("path_ltree", category.PathLtree).Error
	}

	// parent
	var parent models.PostCategory
	if err := r.db.First(&parent, *category.PostCategoryID).Error; err != nil {
		return fmt.Errorf("parent not found: %w", err)
	}

	if parent.PathLtree == "" {
		return fmt.Errorf("parent path_ltree is empty (id=%d)", parent.ID)
	}

	if err := r.db.Create(category).Error; err != nil {
		return err
	}

	category.PathLtree = fmt.Sprintf("%s.%d", parent.PathLtree, category.ID)

	return r.db.Model(category).Update("path_ltree", category.PathLtree).Error
}

func (r *categoryRepository) GetByID(id uint) (*models.PostCategory, error) {
	var category models.PostCategory
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) FindByParam(field string, value any) (*models.PostCategory, error) {
	var category models.PostCategory
	condition := map[string]any{
		field: value,
	}
	if err := r.db.Where(condition).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetDescendants(category *models.PostCategory) ([]*models.PostCategory, error) {
	if category == nil || category.ID == 0 || category.PathLtree == "" {
		return []*models.PostCategory{}, nil
	}

	var descendants []*models.PostCategory
	err := r.db.
		Where("path_ltree <@ ?::ltree AND id <> ?", category.PathLtree, category.ID).
		Order("id ASC").
		Find(&descendants).Error
	if err != nil {
		return nil, err
	}
	return descendants, nil
}

func (r *categoryRepository) GetAllForParent(parent *models.PostCategory) ([]*models.PostCategory, error) {
	if parent == nil || parent.ID == 0 || parent.PathLtree == "" {
		return []*models.PostCategory{}, nil
	}

	var items []*models.PostCategory
	err := r.db.
		Where("NOT (path_ltree <@ ?::ltree) AND id <> ?", parent.PathLtree, parent.ID).
		Order("id ASC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *categoryRepository) Update(new *models.PostCategory, _ *models.PostCategory) error {
	if new == nil {
		return fmt.Errorf("new is nil")
	}

	new.Updating()

	return r.db.Transaction(func(db *gorm.DB) error {
		// 1) читаем old под блокировкой
		var old models.PostCategory
		if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&old, new.ID).Error; err != nil {
			return err
		}

		if err := validatePathEndsWithID(old.PathLtree, old.ID); err != nil {
			return fmt.Errorf("refuse to update because old path_ltree is invalid: %w", err)
		}

		// 2) определяем oldParent/newParent
		oldParent, newParent := uint(0), uint(0)
		if old.PostCategoryID != nil {
			oldParent = *old.PostCategoryID
		}
		if new.PostCategoryID != nil {
			newParent = *new.PostCategoryID
		}

		// 2.0) защита: parent=self
		if newParent != 0 && newParent == old.ID {
			return fmt.Errorf("cannot set parent to self: id=%d", old.ID)
		}

		// (опционально, но рекомендую оставить)
		// чтобы случайно не занулить pointer-поля при Save(new),
		// копируем из old как в MenuItem new.MenuID = old.MenuID
		new.TemplateName = old.TemplateName
		new.UserID = old.UserID

		// 2.1) если родитель не менялся — обычный save, path_ltree оставляем прежним
		if oldParent == newParent {
			new.PathLtree = old.PathLtree
			return db.Select(new.Fields()).Save(new).Error
		}

		// 3) валидируем нового родителя (если перенос не в root)
		var newParentCategory models.PostCategory
		if newParent != 0 {
			if err := db.Select("id", "path_ltree").
				First(&newParentCategory, newParent).Error; err != nil {
				return fmt.Errorf("new parent not found: %w", err)
			}

			if newParentCategory.PathLtree == "" {
				return fmt.Errorf("new parent has empty path_ltree (id=%d)", newParentCategory.ID)
			}

			// защита: нельзя завернуть в своего потомка
			oldPrefix := old.PathLtree
			descPrefix := oldPrefix + "."
			if newParentCategory.PathLtree == oldPrefix || strings.HasPrefix(newParentCategory.PathLtree, descPrefix) {
				return fmt.Errorf(
					"cannot move node under its descendant: node=%d new_parent=%d",
					old.ID, newParentCategory.ID,
				)
			}
		}

		oldPathLtree := old.PathLtree

		var newPathLtree string
		if newParent == 0 {
			newPathLtree = fmt.Sprintf("%d", new.ID)
		} else {
			newPathLtree = fmt.Sprintf("%s.%d", newParentCategory.PathLtree, new.ID)
		}

		// 4) обновляем потомков (кроме корня поддерева) — БЕЗ template_id/user_id
		if err := db.Exec(fmt.Sprintf(`
			UPDATE %s
			SET path_ltree = ?::ltree || subpath(path_ltree, nlevel(?::ltree))
			WHERE path_ltree <@ ?::ltree
			  AND path_ltree <> ?::ltree
		`, new.GetTable()),
			newPathLtree, oldPathLtree,
			oldPathLtree, oldPathLtree,
		).Error; err != nil {
			return fmt.Errorf("failed to update subtree paths: %w", err)
		}

		// 5) обновляем сам узел
		new.PathLtree = newPathLtree
		return db.Select(new.Fields()).Save(new).Error
	})
}

func (r *categoryRepository) UpdateFieldsByUUIDs(uuids []uuid.UUID, patch map[string]any) (int64, error) {
	if len(uuids) == 0 {
		return 0, fmt.Errorf("empty uuids")
	}
	if len(patch) == 0 {
		return 0, fmt.Errorf("empty patch")
	}

	tx := r.db.Model(&models.PostCategory{}).
		Where("uuid IN ?", uuids).
		Updates(patch)

	return tx.RowsAffected, tx.Error
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

func uval(p *uint) uint {
	if p == nil {
		return 0
	}
	return *p
}
