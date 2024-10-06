package models

import (
	"fmt"
	"github.com/axlle-com/blog/pkg/common/db"
	common "github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *Post) error
	GetByID(id uint) (*Post, error)
	Update(post *Post) error
	Delete(id uint) error
	GetAll() ([]*Post, error)
	GetPaginate(paginator contracts.Paginator, filter *PostFilter) ([]*PostFull, error)
	GetByAlias(alias string) (*Post, error)
	GetByAliasNotID(alias string, id uint) (*Post, error)
}

type postRepository struct {
	*common.Paginate
	db *gorm.DB
	c  int
}

func PostRepo() PostRepository {
	return &postRepository{db: db.GetDB()}
}

func (r *postRepository) Create(post *Post) error {
	post.Creating()
	return r.db.Create(post).Error
}

func (r *postRepository) GetByID(id uint) (*Post, error) {
	var model Post
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *postRepository) Update(post *Post) error {
	post.Updating()
	return r.db.Select(
		"UserID",
		"TemplateID",
		"PostCategoryID",
		"MetaTitle",
		"MetaDescription",
		"Alias",
		"URL",
		"IsPublished",
		"IsFavourites",
		"HasComments",
		"ShowImagePost",
		"ShowImageCategory",
		"MakeWatermark",
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
		"Hits",
		"Sort",
		"Stars",
	).Save(post).Error
}

func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&Post{}, id).Error
}

func (r *postRepository) GetAll() ([]*Post, error) {
	var posts []*Post
	if err := r.db.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) GetPaginate(p contracts.Paginator, filter *PostFilter) ([]*PostFull, error) {
	var posts []*PostFull
	var total int64

	query := r.db.Table("posts").
		Select(
			"posts.*",
			"post_categories.title as category_title",
			"post_categories.title_short as category_title_short",
			"templates.title as template_title",
			"templates.name as template_name",
			"users.first_name as user_first_name",
			"users.last_name as user_last_name",
		).
		Joins("left join post_categories on post_categories.id = posts.post_category_id").
		Joins("left join users on users.id = posts.user_id").
		Joins("left join templates on templates.id = posts.template_id")

	// TODO WHERE IN; LIKE
	for col, val := range filter.GetMap() {
		if col == "title" {
			query = query.Where(fmt.Sprintf("posts.%v ilike ?", col), fmt.Sprintf("%%%v%%", val))
			continue
		}
		query = query.Where(fmt.Sprintf("posts.%v = ?", col), val)
	}

	query.Count(&total)

	err := query.Scopes(r.SetPaginate(p.GetPage(), p.GetPageSize())).
		Order("posts.id ASC").
		Scan(&posts).Error

	p.SetTotal(int(total))
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) GetByAlias(alias string) (*Post, error) {
	var post Post
	if err := r.db.Where("alias = ?", alias).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// GetByAliasNotID TODO AND id <> 0 ORDER BY "posts"."id" LIMIT 1
func (r *postRepository) GetByAliasNotID(alias string, id uint) (*Post, error) {
	var post Post
	if err := r.db.Where("alias = ?", alias).Where("id <> ?", id).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}
