package models

import (
	"github.com/axlle-com/blog/app/errors"
	. "github.com/axlle-com/blog/app/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

func NewPostFilter() *PostFilter {
	return &PostFilter{}
}

type PostFilter struct {
	ID             *uint   `json:"id" form:"id" binding:"omitempty"`
	TemplateID     *uint   `json:"template_id" form:"template_id" binding:"omitempty"`
	UserID         *uint   `json:"user_id" form:"user_id" binding:"omitempty"`
	PostCategoryID *uint   `json:"post_category_id" form:"post_category_id" binding:"omitempty"`
	Title          *string `json:"title" form:"title" binding:"omitempty"`
	Date           *string `json:"date" form:"date" binding:"omitempty"`
	Filter
}

func (p *PostFilter) ValidateForm(ctx *gin.Context) (*PostFilter, *errors.Errors) {
	err := p.Filter.ValidateForm(ctx, p)
	return p, err
}

func (p *PostFilter) ValidateQuery(ctx *gin.Context) (*PostFilter, *errors.Errors) {
	err := p.Filter.ValidateQuery(ctx, p)
	return p, err
}

func (p *PostFilter) PrintTemplateID() uint {
	if p.TemplateID == nil {
		return 0
	}
	return *p.TemplateID
}

func (p *PostFilter) PrintUserID() uint {
	if p.UserID == nil {
		return 0
	}
	return *p.UserID
}

func (p *PostFilter) PrintPostCategoryID() uint {
	if p.PostCategoryID == nil {
		return 0
	}
	return *p.PostCategoryID
}

func (p *PostFilter) GetURL() string {
	if p.GetQueryString() == "" {
		return "posts"
	}
	return string("posts?" + p.GetQueryString())
}

func (p *PostFilter) PrintID() string {
	if p.ID == nil {
		return ""
	}
	return strconv.Itoa(int(*p.ID))
}
