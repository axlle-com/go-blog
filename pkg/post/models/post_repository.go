package models

import (
	"github.com/axlle-com/blog/pkg/common/db"
	common "github.com/axlle-com/blog/pkg/common/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *Post) error
	GetByID(id uint) (*Post, error)
	Update(post *Post) error
	Delete(id uint) error
	GetAll() ([]Post, error)
	GetPaginate(page, pageSize int) ([]PostResponse, int, error)
	GetByAlias(alias string) (*Post, error)
	GetByAliasNotID(alias string, id uint) (*Post, error)
}

type postRepository struct {
	*common.Paginate
	db *gorm.DB
	c  int
}

func NewPostRepo() PostRepository {
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

func (r *postRepository) GetAll() ([]Post, error) {
	var posts []Post
	if err := r.db.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) GetPaginate(page, pageSize int) ([]PostResponse, int, error) {
	var posts []PostResponse
	var total int64

	r.db.Model(&Post{}).Count(&total)
	err := r.db.Table("posts").
		Scopes(r.SetPaginate(page, pageSize)).
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
		Joins("left join templates on templates.id = posts.template_id").
		Order("posts.id ASC").
		Scan(&posts).Error
	if err != nil {
		return nil, 0, err
	}
	return posts, int(total), nil
}

func (r *postRepository) GetByAlias(alias string) (*Post, error) {
	var post Post
	if err := r.db.Where("alias = ?", alias).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) GetByAliasNotID(alias string, id uint) (*Post, error) {
	var post Post
	if err := r.db.Where("alias = ?", alias).Where("id <> ?", id).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}
