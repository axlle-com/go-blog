package request

import (
	"github.com/axlle-com/blog/app/models/contract"
	post "github.com/axlle-com/blog/pkg/blog/models"
)

func NewPostResponse() *PostResponse {
	return &PostResponse{}
}

type PostResponse struct {
	Post       *post.Post
	Categories []*post.PostCategory
	Templates  []contract.Template
}
