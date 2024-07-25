package repository

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/models"
	post "github.com/axlle-com/blog/pkg/post/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *models.Post) error
	GetByID(id uint) (*models.Post, error)
	Update(post *models.Post) error
	Delete(id uint) error
	GetAll() ([]models.Post, error)
	GetPaginate(page, pageSize int) ([]post.Post, int, error)
}

type repository struct {
	*models.Paginate
	db *gorm.DB
}

func NewRepository() PostRepository {
	return &repository{db: db.GetDB()}
}

func (r *repository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *repository) GetByID(id uint) (*models.Post, error) {
	var model models.Post
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *repository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}

func (r *repository) GetAll() ([]models.Post, error) {
	var posts []models.Post
	if err := r.db.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *repository) GetPaginate(page, pageSize int) ([]post.Post, int, error) {
	var posts []post.Post
	var total int64

	r.db.Model(&post.Post{}).Count(&total)
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
		return nil, 0, err
	}
	return posts, int(total), nil
}
