package request

import (
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/gin-gonic/gin"
)

type Response struct {
	paginator *contract.Paginator
}

func NewResponse(p *contract.Paginator) *Response {
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
