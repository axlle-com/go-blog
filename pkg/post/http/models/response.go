package models

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	common "github.com/axlle-com/blog/pkg/common/models"
	. "github.com/axlle-com/blog/pkg/post/models"
	template "github.com/axlle-com/blog/pkg/template/repository"
	userRepo "github.com/axlle-com/blog/pkg/user/repository"
	"github.com/gin-gonic/gin"
)

type Response struct {
	posts      []*PostResponse
	users      []*common.User
	categories []*PostCategory
	templates  []*common.Template
	paginator  common.Paginator
}

func NewResponse(p common.Paginator) *Response {
	return &Response{paginator: p}
}

func (r *Response) GetForWeb() gin.H {
	users, err := userRepo.NewRepo().GetAll()
	if err != nil {
		logger.New().Error(err)
	}

	categories, err := NewCategoryRepo().GetAll()
	if err != nil {
		logger.New().Error(err)
	}

	templates, err := template.NewRepo().GetAllTemplates()
	if err != nil {
		logger.New().Error(err)
	}

	posts, total, err := NewPostRepo().GetPaginate(r.paginator.GetPage(), r.paginator.GetPageSize())
	if err != nil {
		logger.New().Error(err)
	}

	r.paginator.SetTotal(total)

	return gin.H{
		"title":      "Страница постов",
		"posts":      posts,
		"categories": categories,
		"templates":  templates,
		"users":      users,
		"total":      total,
		"paginator":  r.paginator,
	}
}

func (r *Response) GetForAjax() gin.H {
	users, err := userRepo.NewRepo().GetAll()
	if err != nil {
		logger.New().Error(err)
	}

	categories, err := NewCategoryRepo().GetAll()
	if err != nil {
		logger.New().Error(err)
	}

	templates, err := template.NewRepo().GetAllTemplates()
	if err != nil {
		logger.New().Error(err)
	}

	posts, total, err := NewPostRepo().GetPaginate(r.paginator.GetPage(), r.paginator.GetPageSize())
	if err != nil {
		logger.New().Error(err)
	}

	r.paginator.SetTotal(total)

	return gin.H{
		"title":      "Страница постов",
		"posts":      posts,
		"categories": categories,
		"templates":  templates,
		"users":      users,
		"total":      total,
		"paginator":  r.paginator,
	}
}
