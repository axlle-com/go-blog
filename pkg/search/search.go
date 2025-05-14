package search

import "github.com/axlle-com/blog/pkg/blog/models"

type Search interface {
	CreateIndex(string) error
	AddPost(*models.Post) error
	AddPostCategory(*models.PostCategory) error
	SearchPosts(string) ([]models.Post, error)
}
