package repository

import (
	"errors"
	"fmt"
	"strings"

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
}

type infoBlockRepository struct {
	db *gorm.DB
	*app.Paginate
}

func NewInfoBlockRepo(db *gorm.DB) InfoBlockRepository {
	r := &infoBlockRepository{db: db}
	return r
}

func (r *infoBlockRepository) WithTx(tx *gorm.DB) InfoBlockRepository {
	return &infoBlockRepository{db: tx}
}

func (r *infoBlockRepository) Create(infoBlock *models.InfoBlock) error {
	infoBlock.Creating()
	if infoBlock.InfoBlockID == nil || *infoBlock.InfoBlockID == 0 {
		if err := r.db.Create(infoBlock).Error; err != nil {
			return err
		}
		// Для корневого инфоблока путь – просто /id/
		infoBlock.Path = fmt.Sprintf("/%d/", infoBlock.ID)
		return r.db.Model(infoBlock).Update("path", infoBlock.Path).Error
	}

	// Если есть родитель, получаем его данные.
	var parent models.InfoBlock
	if err := r.db.First(&parent, *infoBlock.InfoBlockID).Error; err != nil {
		return fmt.Errorf("parent not found: %w", err)
	}

	if err := r.db.Create(infoBlock).Error; err != nil {
		return err
	}
	// Путь дочернего узла – путь родителя + id дочернего.
	infoBlock.Path = fmt.Sprintf("%s%d/", parent.Path, infoBlock.ID)
	return r.db.Model(infoBlock).Update("path", infoBlock.Path).Error
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

	// Применяем фильтры напрямую из полей структуры
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

func (r *infoBlockRepository) save(infoBlock *models.InfoBlock) error {
	return r.db.Select(
		"TemplateID",
		"InfoBlockID",
		"Media",
		"Title",
		"Description",
		"Image",
		"Path",
	).Save(infoBlock).Error
}

func (r *infoBlockRepository) Update(new *models.InfoBlock, old *models.InfoBlock) error {
	new.Updating()
	new.Path = old.Path

	// Если родитель не изменился – просто сохраняем изменения.
	oldParent, newParent := uint(0), uint(0)
	if old.InfoBlockID != nil {
		oldParent = *old.InfoBlockID
	}
	if new.InfoBlockID != nil {
		newParent = *new.InfoBlockID
	}

	if oldParent == newParent {
		return r.save(new)
	}

	// Если родитель меняется, требуется пересчитать путь для нового поддерева.
	var newParentInfoBlock models.InfoBlock
	if newParent != 0 {
		if err := r.db.First(&newParentInfoBlock, newParent).Error; err != nil {
			return fmt.Errorf("new parent not found: %w", err)
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
		newPath = fmt.Sprintf("%s%d/", newParentInfoBlock.Path, new.ID)
	}

	// Обновляем путь для узла и всех потомков.
	var descendants []*models.InfoBlock
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

func (r *infoBlockRepository) Delete(infoBlock *models.InfoBlock) error {
	if !infoBlock.Deleting() {
		return errors.New("deletion errors occurred")
	}

	likePattern := fmt.Sprintf("%s%%", infoBlock.Path)
	return r.db.Where("path LIKE ?", likePattern).Delete(&models.InfoBlock{}).Error
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
	query := r.db.Where("id IN ?", ids)

	if err := query.Find(&infoBlocks).Error; err != nil {
		return nil, err
	}

	return infoBlocks, nil
}

func (r *infoBlockRepository) DeleteByIDs(ids []uint) (err error) {
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
	likePattern := fmt.Sprintf("%s%%", infoBlock.Path)
	err := r.db.
		Where("path LIKE ? AND id <> ?", likePattern, infoBlock.ID).
		Order("id ASC").
		Find(&descendants).Error
	if err != nil {
		return nil, err
	}

	return descendants, err
}

func (r *infoBlockRepository) GetAllForParent(parent *models.InfoBlock) ([]*models.InfoBlock, error) {
	var items []*models.InfoBlock
	likePattern := fmt.Sprintf("%s%%", parent.Path)
	err := r.db.
		Where("path NOT LIKE ? AND id <> ?", likePattern, parent.ID).
		Order("id ASC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, err
}
