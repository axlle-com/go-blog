package repository

import (
	"github.com/axlle-com/blog/pkg/common/db"
	common "github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/post/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *models.Post) error
	GetByID(id uint) (*models.Post, error)
	Update(post *models.Post) error
	Delete(id uint) error
	GetAll() ([]models.Post, error)
	GetPaginate(page, pageSize int) ([]models.PostResponse, int, error)
}

type postRepository struct {
	*common.Paginate
	db *gorm.DB
}

func NewPostRepository() PostRepository {
	return &postRepository{db: db.GetDB()}
}

func (r *postRepository) Create(post *models.Post) error {
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
	return r.db.Delete(&models.Post{}, id).Error
}

func (r *postRepository) GetAll() ([]models.Post, error) {
	var posts []models.Post
	if err := r.db.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) GetPaginate(page, pageSize int) ([]models.PostResponse, int, error) {
	var posts []models.PostResponse
	var total int64

	r.db.Model(&models.Post{}).Count(&total)
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
