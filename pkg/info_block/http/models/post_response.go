package models

import (
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	post "github.com/axlle-com/blog/pkg/post/models"
)

func NewPostResponse() *PostResponse {
	return &PostResponse{}
}

type PostResponse struct {
	Post       *post.Post
	Categories []*post.PostCategory
	Templates  []contracts.Template
}
