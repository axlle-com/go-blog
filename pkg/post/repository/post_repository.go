package repository

import (
	"github.com/axlle-com/blog/pkg/common/db"
	models "github.com/axlle-com/blog/pkg/common/models"
	post "github.com/axlle-com/blog/pkg/post/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreatePost(post *models.Post) error
	GetPostByID(id uint) (*models.Post, error)
	UpdatePost(post *models.Post) error
	DeletePost(id uint) error
	GetAllPosts() ([]models.Post, error)
	GetPaginate(page, pageSize int) ([]post.Post, error)
}

type repository struct {
	*models.Paginate
	db *gorm.DB
}

func NewRepository() Repository {
	return &repository{db: db.GetDB()}
}

func (r *repository) CreatePost(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *repository) GetPostByID(id uint) (*models.Post, error) {
	var model models.Post
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *repository) UpdatePost(post *models.Post) error {
	return r.db.Save(post).Error
}

func (r *repository) DeletePost(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}

func (r *repository) GetAllPosts() ([]models.Post, error) {
	var posts []models.Post
	if err := r.db.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *repository) GetPaginate(page, pageSize int) ([]post.Post, error) {
	var posts []post.Post

	err := r.db.Table("posts").
		Scopes(r.SetPaginate(page, pageSize)).
		Select("posts.*," +
			"post_categories.title as category_title," +
			"post_categories.title_short as category_title_short," +
			"templates.title as template_title," +
			"templates.name as template_name," +
			"users.first_name as user_first_name," +
			"users.last_name as user_last_name").
		Joins("left join post_categories on post_categories.id = posts.post_category_id").
		Joins("left join users on users.id = posts.user_id").
		Joins("left join templates on templates.id = posts.template_id").
		Scan(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}
