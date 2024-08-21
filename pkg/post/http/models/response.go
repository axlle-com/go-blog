package models

import (
	common "github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"github.com/axlle-com/blog/pkg/post/models"
	"github.com/gin-gonic/gin"
)

type Response struct {
	posts      []*models.PostResponse
	users      []*common.User
	categories []*models.PostCategory
	templates  []contracts.Template
	paginator  common.Paginator
}

func NewResponse(p common.Paginator) *Response {
	return &Response{paginator: p}
}

func (r *Response) GetForWeb() gin.H {

	return gin.H{
		"title":     "Страница постов",
		"paginator": r.paginator,
	}
}

func (r *Response) GetForAjax() gin.H {

	return gin.H{
		"title":     "Страница постов",
		"paginator": r.paginator,
	}
}
