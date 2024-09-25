package models

import (
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"github.com/gin-gonic/gin"
)

type Response struct {
	paginator *contracts.Paginator
}

func NewResponse(p *contracts.Paginator) *Response {
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
