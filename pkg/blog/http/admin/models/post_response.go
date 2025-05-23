package models

import (
	"github.com/axlle-com/blog/app/models/contracts"
	post "github.com/axlle-com/blog/pkg/blog/models"
)

func NewPostResponse() *PostResponse {
	return &PostResponse{}
}

type PostResponse struct {
	Post       *post.Post
	Categories []*post.PostCategory
	Templates  []contracts.Template
}
