package repository

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	CreatePost(post *models.Post) error
	GetPostByID(id uint) (*models.Post, error)
	GetPostByEmail(email string) (*models.Post, error)
	UpdatePost(post *models.Post) error
	DeletePost(id uint) error
	GetAllPosts() ([]models.Post, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository() PostRepository {
	return &postRepository{db: db.GetDB()}
}

func (r *postRepository) CreatePost(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) GetPostByID(id uint) (*models.Post, error) {
	var post models.Post
	if err := r.db.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) GetPostByEmail(email string) (*models.Post, error) {
	var post models.Post
	if err := r.db.Where("email = ?", email).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) UpdatePost(post *models.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) DeletePost(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}

func (r *postRepository) GetAllPosts() ([]models.Post, error) {
	var posts []models.Post
	if err := r.db.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
