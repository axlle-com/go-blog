package repository

import (
	"fmt"

	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostTagRepository interface {
	Create(postTag *models.PostTag) error
	GetByID(id uint) (*models.PostTag, error)
	GetByIDs(ids []uint) ([]*models.PostTag, error)
	GetByNames(titles []string) ([]*models.PostTag, error)
	Update(postTag *models.PostTag) error
	UpdateFieldsByUUIDs(uuids []uuid.UUID, patch map[string]any) (int64, error)
	DeleteByID(id uint) error
	Delete(*models.PostTag) error
	DeleteByIDs(ids []uint) (err error)
	GetAll() ([]*models.PostTag, error)
	GetAllIds() ([]uint, error)
	GetForResource(contract.Resource) ([]*models.PostTag, error)
	WithPaginate(p contract.Paginator, filter *models.TagFilter) ([]*models.PostTag, error)
	WithTx(tx *gorm.DB) PostTagRepository
}

type postTagRepository struct {
	db *gorm.DB
	*app.Paginate
}

func NewPostTagRepo(db *gorm.DB) PostTagRepository {
	r := &postTagRepository{db: db}
	return r
}

func (r *postTagRepository) WithTx(tx *gorm.DB) PostTagRepository {
	return &postTagRepository{db: tx}
}

func (r *postTagRepository) WithPaginate(p contract.Paginator, filter *models.TagFilter) ([]*models.PostTag, error) {
	var postTags []*models.PostTag
	var total int64

	tag := models.PostTag{}
	table := tag.GetTable()

	query := r.db.Model(&tag)

	if filter != nil {
		if filter.ID != nil {
			query = query.Where(fmt.Sprintf("%s.id = ?", table), *filter.ID)
		}
		if filter.TemplateID != nil {
			query = query.Where(fmt.Sprintf("%s.template_id = ?", table), *filter.TemplateID)
		}
		if filter.Name != nil && *filter.Name != "" {
			query = query.Where(fmt.Sprintf("%s.name = ?", table), *filter.Name)
		}
		if filter.Title != nil && *filter.Title != "" {
			query = query.Where(fmt.Sprintf("%s.title ilike ?", table), fmt.Sprintf("%%%s%%", *filter.Title))
		}
		if filter.Query != nil && *filter.Query != "" {
			query = query.Where(fmt.Sprintf("(%s.title ilike ? OR %s.name ilike ?)", table, table), fmt.Sprintf("%%%s%%", *filter.Query), fmt.Sprintf("%%%s%%", *filter.Query))
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
		Find(&postTags).Error
	if err != nil {
		return nil, err
	}

	p.SetTotal(int(total))
	return postTags, nil
}

func (r *postTagRepository) Create(postTag *models.PostTag) error {
	postTag.Creating()
	return r.db.Create(postTag).Error
}

func (r *postTagRepository) GetByID(id uint) (*models.PostTag, error) {
	var postTag models.PostTag
	if err := r.db.First(&postTag, id).Error; err != nil {
		return nil, err
	}
	return &postTag, nil
}

func (r *postTagRepository) GetByIDs(ids []uint) ([]*models.PostTag, error) {
	var galleries []*models.PostTag

	if err := r.db.Where("id IN ?", ids).Find(&galleries).Error; err != nil {
		return nil, err
	}
	return galleries, nil
}

func (r *postTagRepository) GetByNames(names []string) ([]*models.PostTag, error) {
	var tags []*models.PostTag

	if err := r.db.Where("name IN ?", names).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *postTagRepository) Update(postTag *models.PostTag) error {
	postTag.Updating()
	return r.db.
		Select(
			"TemplateID",
			"Name",
			"Title",
			"Description",
			"Image",
			"MetaTitle",
			"MetaDescription",
			"Alias",
			"URL",
		).
		Save(postTag).Error
}

func (r *postTagRepository) DeleteByID(id uint) error {
	return r.db.Delete(models.PostTag{}, id).Error
}

func (r *postTagRepository) Delete(g *models.PostTag) (err error) {
	return r.db.Delete(g, g.ID).Error
}

func (r *postTagRepository) DeleteByIDs(ids []uint) (err error) {
	return r.db.Where("id IN ?", ids).Delete(&models.PostTag{}).Error
}

func (r *postTagRepository) GetAll() ([]*models.PostTag, error) {
	var galleries []*models.PostTag
	if err := r.db.Order("id ASC").Find(&galleries).Error; err != nil {
		return nil, err
	}
	return galleries, nil
}

func (r *postTagRepository) GetForResource(resource contract.Resource) ([]*models.PostTag, error) {
	var galleries []*models.PostTag
	query := r.db.
		Joins("inner join post_tag_has_resources as r on post_tags.id = r.post_tag_id").
		Where("r.resource_uuid = ?", resource.GetUUID()).
		Model(&models.PostTag{})

	err := query.Find(&galleries).Error
	return galleries, err
}

func (r *postTagRepository) GetAllIds() ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&models.PostTag{}).Pluck("id", &ids).Error; err != nil {
		logger.Error(err)
	}
	return ids, nil
}

func (r *postTagRepository) UpdateFieldsByUUIDs(uuids []uuid.UUID, patch map[string]any) (int64, error) {
	if len(uuids) == 0 {
		return 0, fmt.Errorf("empty uuids")
	}
	if len(patch) == 0 {
		return 0, fmt.Errorf("empty patch")
	}

	tx := r.db.Model(&models.PostTag{}).
		Where("uuid IN ?", uuids).
		Updates(patch)

	return tx.RowsAffected, tx.Error
}
