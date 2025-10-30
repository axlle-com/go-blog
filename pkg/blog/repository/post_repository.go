package repository

import (
	"errors"
	"fmt"

	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepository interface {
	WithTx(tx *gorm.DB) PostRepository
	Create(post *models.Post) error
	GetByID(id uint) (*models.Post, error)
	GetByParam(field string, value any) (*models.Post, error)
	GetByParams(params map[string]any) ([]*models.Post, error)
	Update(post *models.Post) error
	UpdateFieldsByUUIDs(uuids []uuid.UUID, patch map[string]any) (int64, error)
	Delete(post *models.Post) error
	GetAll() ([]*models.Post, error)
	WithPaginate(paginator contract.Paginator, filter *models.PostFilter) ([]*models.Post, error)
	GetByAlias(alias string) (*models.Post, error)
	GetByAliasNotID(alias string, id uint) (*models.Post, error)
}

type postRepository struct {
	db *gorm.DB
	*app.Paginate
}

func NewPostRepo(db *gorm.DB) PostRepository {
	r := &postRepository{db: db}
	return r
}

func (r *postRepository) WithTx(tx *gorm.DB) PostRepository {
	return &postRepository{db: tx}
}

func (r *postRepository) Create(post *models.Post) error {
	post.Creating()
	return r.db.Create(post).Error
}

func (r *postRepository) GetByID(id uint) (*models.Post, error) {
	var model models.Post
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *postRepository) Update(post *models.Post) error {
	post.Updating()
	return r.db.Select(
		"UserID",
		"TemplateID",
		"PostCategoryID",
		"MetaTitle",
		"MetaDescription",
		"Alias",
		"URL",
		"IsMain",
		"IsPublished",
		"IsFavourites",
		"HasComments",
		"ShowImagePost",
		"ShowImageCategory",
		"InSitemap",
		"Media",
		"Title",
		"TitleShort",
		"DescriptionPreview",
		"Description",
		"ShowDate",
		"DatePub",
		"DateEnd",
		"Image",
		"Sort",
	).Save(post).Error
}

func (r *postRepository) Delete(post *models.Post) error {
	if post.Deleting() {
		return r.db.Delete(&models.Post{}, post.ID).Error
	}
	return errors.New("при удалении произошли ошибки")
}

func (r *postRepository) GetAll() ([]*models.Post, error) {
	var posts []*models.Post
	if err := r.db.Order("id ASC").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) WithPaginate(p contract.Paginator, filter *models.PostFilter) ([]*models.Post, error) {
	var posts []*models.Post
	var total int64

	post := models.Post{}
	table := post.GetTable()

	query := r.db.Model(&posts)

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
		Scan(&posts).Error

	p.SetTotal(int(total))
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) GetByParam(field string, value any) (*models.Post, error) {
	var post models.Post
	condition := map[string]any{
		field: value,
	}
	if err := r.db.Where(condition).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) GetByParams(params map[string]any) ([]*models.Post, error) {
	var posts []*models.Post
	if err := r.db.Where(params).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) GetByAlias(alias string) (*models.Post, error) {
	var post models.Post
	if err := r.db.Where("alias = ?", alias).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// GetByAliasNotID TODO AND id <> 0 ORDER BY "posts"."id" LIMIT 1
func (r *postRepository) GetByAliasNotID(alias string, id uint) (*models.Post, error) {
	var post models.Post
	if err := r.db.Where("alias = ?", alias).Where("id <> ?", id).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) UpdateFieldsByUUIDs(uuids []uuid.UUID, patch map[string]any) (int64, error) {
	if len(uuids) == 0 {
		return 0, fmt.Errorf("empty uuids")
	}
	if len(patch) == 0 {
		return 0, fmt.Errorf("empty patch")
	}

	tx := r.db.Model(&models.Post{}).
		Where("uuid IN ?", uuids).
		Updates(patch)

	return tx.RowsAffected, tx.Error
}
