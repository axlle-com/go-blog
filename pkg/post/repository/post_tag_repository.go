package repository

import (
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/post/models"
	"gorm.io/gorm"
)

type PostTagRepository interface {
	Create(postTag *models.PostTag) error
	GetByID(id uint) (*models.PostTag, error)
	GetByIDs(ids []uint) ([]*models.PostTag, error)
	Update(postTag *models.PostTag) error
	DeleteByID(id uint) error
	Delete(*models.PostTag) error
	DeleteByIDs(ids []uint) (err error)
	GetAll() ([]*models.PostTag, error)
	GetAllIds() ([]uint, error)
	GetForResource(contracts.Resource) ([]*models.PostTag, error)
	WithTx(tx *gorm.DB) PostTagRepository
}

type postTagRepository struct {
	db *gorm.DB
	*app.Paginate
	withImages bool
}

func NewPostTagRepo(db contracts.DB) PostTagRepository {
	r := &postTagRepository{db: db.GORM()}
	return r
}

func (r *postTagRepository) WithTx(tx *gorm.DB) PostTagRepository {
	return &postTagRepository{db: tx}
}

func (r *postTagRepository) Create(postTag *models.PostTag) error {
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

func (r *postTagRepository) Update(postTag *models.PostTag) error {
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

func (r *postTagRepository) GetForResource(resource contracts.Resource) ([]*models.PostTag, error) {
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
