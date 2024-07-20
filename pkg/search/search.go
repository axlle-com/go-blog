package search

import "github.com/axlle-com/blog/pkg/common/models"

type Search interface {
	CreateIndex(string) error
	AddPost(*models.Post) error
	SearchPosts(string) ([]models.Post, error)
}
