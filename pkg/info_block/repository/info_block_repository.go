package repository

import (
	"errors"
	"fmt"

	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"gorm.io/gorm"
)

type InfoBlockRepository interface {
	WithTx(tx *gorm.DB) InfoBlockRepository
	Create(infoBlock *models.InfoBlock) error
	Update(new *models.InfoBlock, old *models.InfoBlock) error
	GetAll() ([]*models.InfoBlock, error)
	GetByIDs(ids []uint) ([]*models.InfoBlock, error)
	GetForResourceByFilter(filter *models.InfoBlockFilter) ([]*models.InfoBlockResponse, error)

	FindByID(id uint) (*models.InfoBlock, error)
	FindByFilter(filter *models.InfoBlockFilter) (*models.InfoBlock, error)

	WithPaginate(p contract.Paginator, filter *models.InfoBlockFilter) ([]*models.InfoBlock, error)

	Delete(infoBlock *models.InfoBlock) error
	DeleteByIDs(ids []uint) (err error)

	GetRoots() ([]*models.InfoBlock, error)
	GetDescendants(infoBlock *models.InfoBlock) ([]*models.InfoBlock, error)
	GetAllForParent(parent *models.InfoBlock) ([]*models.InfoBlock, error)

	// Потомки для нескольких корней
	GetDescendantsByRoots(rootIDs []uint) ([]*models.InfoBlock, error)
}

type infoBlockRepository struct {
	db *gorm.DB
	*app.Paginate
}

func NewInfoBlockRepo(db *gorm.DB) InfoBlockRepository {
	return &infoBlockRepository{db: db}
}

func (r *infoBlockRepository) WithTx(tx *gorm.DB) InfoBlockRepository {
	return &infoBlockRepository{db: tx}
}

func (r *infoBlockRepository) Create(infoBlock *models.InfoBlock) error {
	infoBlock.Creating()

	// root
	if infoBlock.InfoBlockID == nil || *infoBlock.InfoBlockID == 0 {
		if err := r.db.Create(infoBlock).Error; err != nil {
			return err
		}

		infoBlock.PathLtree = fmt.Sprintf("%d", infoBlock.ID)

		return r.db.Model(infoBlock).UpdateColumn("path_ltree", infoBlock.PathLtree).Error
	}

	// parent
	var parent models.InfoBlock
	if err := r.db.First(&parent, *infoBlock.InfoBlockID).Error; err != nil {
		return fmt.Errorf("parent not found: %w", err)
	}

	if parent.PathLtree == "" {
		return fmt.Errorf("parent path_ltree is empty (id=%d)", parent.ID)
	}

	if err := r.db.Create(infoBlock).Error; err != nil {
		return err
	}

	infoBlock.PathLtree = fmt.Sprintf("%s.%d", parent.PathLtree, infoBlock.ID)

	return r.db.Model(infoBlock).UpdateColumn("path_ltree", infoBlock.PathLtree).Error
}

func (r *infoBlockRepository) FindByID(id uint) (*models.InfoBlock, error) {
	var model models.InfoBlock
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *infoBlockRepository) FindByFilter(filter *models.InfoBlockFilter) (*models.InfoBlock, error) {
	var model models.InfoBlock
	query := r.db.Model(&model)

	if filter.Title != nil {
		query = query.Where("title = ?", *filter.Title)
	}

	if filter.TemplateID != nil {
		query = query.Where("template_id = ?", *filter.TemplateID)
	}

	if filter.ID != nil {
		query = query.Where("id = ?", *filter.ID)
	}

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	if err := query.First(&model).Error; err != nil {
		return nil, err
	}

	return &model, nil
}

func (r *infoBlockRepository) WithPaginate(p contract.Paginator, filter *models.InfoBlockFilter) ([]*models.InfoBlock, error) {
	var infoBlocks []*models.InfoBlock
	var total int64

	infoBlock := models.InfoBlock{}
	table := infoBlock.GetTable()

	query := r.db.Model(&infoBlock)

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
		Order(fmt.Sprintf("%s.id ASC", infoBlock.GetTable())).
		Find(&infoBlocks).Error
	if err != nil {
		return nil, err
	}

	p.SetTotal(int(total))
	return infoBlocks, nil
}

func (r *infoBlockRepository) Update(new *models.InfoBlock, old *models.InfoBlock) error {
	new.Updating()

	oldParent, newParent := uint(0), uint(0)
	if old.InfoBlockID != nil {
		oldParent = *old.InfoBlockID
	}
	if new.InfoBlockID != nil {
		newParent = *new.InfoBlockID
	}

	// parent not changed => обычное сохранение (path_ltree не трогаем)
	if oldParent == newParent {
		return r.db.Select(
			"TemplateID",
			"InfoBlockID",
			"Media",
			"Title",
			"Description",
			"Image",
		).Save(new).Error
	}

	if old.PathLtree == "" {
		return fmt.Errorf("old path_ltree is empty (id=%d)", old.ID)
	}

	// получаем нового родителя (если есть)
	var newParentInfoBlock models.InfoBlock
	if newParent != 0 {
		if err := r.db.First(&newParentInfoBlock, newParent).Error; err != nil {
			return fmt.Errorf("new parent not found: %w", err)
		}
		if newParentInfoBlock.PathLtree == "" {
			return fmt.Errorf("new parent path_ltree is empty (id=%d)", newParentInfoBlock.ID)
		}
	}

	// новый путь для корня поддерева
	var newPathLtree string
	if newParent == 0 {
		newPathLtree = fmt.Sprintf("%d", new.ID)
	} else {
		newPathLtree = fmt.Sprintf("%s.%d", newParentInfoBlock.PathLtree, new.ID)
	}

	tx := r.db.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// 1) одним запросом обновляем path_ltree у узла и всех потомков
	// path_ltree = newPath || subpath(path_ltree, nlevel(oldPath))
	// WHERE path_ltree <@ oldPath
	if err := tx.Model(&models.InfoBlock{}).
		Where("path_ltree <@ ?::ltree", old.PathLtree).
		UpdateColumn(
			"path_ltree",
			gorm.Expr("?::ltree || subpath(path_ltree, nlevel(?::ltree))", newPathLtree, old.PathLtree),
		).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update subtree paths: %w", err)
	}

	// 2) сохраняем поля узла (включая нового родителя и новый path_ltree)
	new.PathLtree = newPathLtree
	if err := tx.Select(
		"TemplateID",
		"InfoBlockID",
		"Media",
		"Title",
		"Description",
		"Image",
		"PathLtree",
	).Save(new).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *infoBlockRepository) Delete(infoBlock *models.InfoBlock) error {
	if !infoBlock.Deleting() {
		return errors.New("deletion errors occurred")
	}
	if infoBlock.PathLtree == "" {
		return fmt.Errorf("path_ltree is empty (id=%d)", infoBlock.ID)
	}

	// узел + потомки
	return r.db.Where("path_ltree <@ ?::ltree", infoBlock.PathLtree).
		Delete(&models.InfoBlock{}).Error
}

func (r *infoBlockRepository) GetAll() ([]*models.InfoBlock, error) {
	var infoBlocks []*models.InfoBlock
	if err := r.db.Order("id ASC").Find(&infoBlocks).Error; err != nil {
		return nil, err
	}

	return infoBlocks, nil
}

func (r *infoBlockRepository) GetByIDs(ids []uint) ([]*models.InfoBlock, error) {
	var infoBlocks []*models.InfoBlock
	if len(ids) == 0 {
		return []*models.InfoBlock{}, nil
	}

	if err := r.db.Where("id IN ?", ids).Find(&infoBlocks).Error; err != nil {
		return nil, err
	}
	return infoBlocks, nil
}

func (r *infoBlockRepository) DeleteByIDs(ids []uint) (err error) {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Where("id IN ?", ids).Delete(&models.InfoBlock{}).Error
}

func (r *infoBlockRepository) GetForResourceByFilter(filter *models.InfoBlockFilter) ([]*models.InfoBlockResponse, error) {
	var infoBlocks []*models.InfoBlockResponse
	query := r.db.
		Joins("inner join info_block_has_resources as r on info_blocks.id = r.info_block_id").
		Select("info_blocks.*", "r.id as relation_id", "r.sort as sort", "r.position as position", "r.resource_uuid as resource_uuid")

	if filter.ResourceUUID != nil {
		query = query.Where("r.resource_uuid = ?", filter.ResourceUUID)
	}

	if filter.RelationID != nil {
		query = query.Where("r.id = ?", filter.RelationID)
	}

	if filter.RelationIDs != nil && len(filter.RelationIDs) > 0 {
		query = query.Where("r.id IN ?", filter.RelationIDs)
	}

	if filter.ID != nil {
		query = query.Where("info_blocks.id = ?", filter.ID)
	}

	if filter.IDs != nil && len(filter.IDs) > 0 {
		query = query.Where("info_blocks.id IN ?", filter.IDs)
	}

	if filter.UUIDs != nil && len(filter.UUIDs) > 0 {
		query = query.Where("info_blocks.uuid IN ?", filter.UUIDs)
	}

	query = query.Order("r.sort ASC").
		Order("r.position ASC").
		Order("r.id DESC").
		Model(&models.InfoBlock{})

	err := query.Find(&infoBlocks).Error
	return infoBlocks, err
}

func (r *infoBlockRepository) GetRoots() ([]*models.InfoBlock, error) {
	var roots []*models.InfoBlock
	err := r.db.Where("info_block_id IS NULL").Find(&roots).Error
	return roots, err
}

func (r *infoBlockRepository) GetDescendants(infoBlock *models.InfoBlock) ([]*models.InfoBlock, error) {
	var descendants []*models.InfoBlock
	if infoBlock == nil || infoBlock.ID == 0 || infoBlock.PathLtree == "" {
		return []*models.InfoBlock{}, nil
	}

	err := r.db.
		Where("path_ltree <@ ?::ltree AND id <> ?", infoBlock.PathLtree, infoBlock.ID).
		Order("id ASC").
		Find(&descendants).Error

	return descendants, err
}

func (r *infoBlockRepository) GetAllForParent(parent *models.InfoBlock) ([]*models.InfoBlock, error) {
	// Сохраняю семантику как в твоём старом коде:
	// "все элементы, которые НЕ являются потомками parent (и не сам parent)"
	var items []*models.InfoBlock
	if parent == nil || parent.ID == 0 || parent.PathLtree == "" {
		return []*models.InfoBlock{}, nil
	}

	err := r.db.
		Where("NOT (path_ltree <@ ?::ltree) AND id <> ?", parent.PathLtree, parent.ID).
		Order("id ASC").
		Find(&items).Error

	return items, err
}

func (r *infoBlockRepository) GetDescendantsByRoots(rootIDs []uint) ([]*models.InfoBlock, error) {
	if len(rootIDs) == 0 {
		return []*models.InfoBlock{}, nil
	}

	// берём пути корней
	type row struct {
		ID        uint
		PathLtree string
	}
	var roots []row
	if err := r.db.Model(&models.InfoBlock{}).
		Select("id, path_ltree").
		Where("id IN ?", rootIDs).
		Find(&roots).Error; err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(roots))
	for _, rt := range roots {
		if rt.PathLtree != "" {
			paths = append(paths, rt.PathLtree)
		}
	}
	if len(paths) == 0 {
		return []*models.InfoBlock{}, nil
	}

	// OR-условия без Raw/unnest
	q := r.db.Model(&models.InfoBlock{}).Where("1=0")
	for _, p := range paths {
		q = q.Or("path_ltree <@ ?::ltree", p)
	}

	var descendants []*models.InfoBlock
	err := q.
		Where("id NOT IN ?", rootIDs).
		Order("id ASC").
		Find(&descendants).Error

	return descendants, err
}
