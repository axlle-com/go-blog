package models

import (
	"strconv"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewCategoryFilterFilter() *CategoryFilter {
	return &CategoryFilter{}
}

type CategoryFilter struct {
	ID             *uint       `json:"id" form:"id" binding:"omitempty"`
	TemplateID     *uint       `json:"template_id" form:"template_id" binding:"omitempty"`
	UserID         *uint       `json:"user_id" form:"user_id" binding:"omitempty"`
	PostCategoryID *uint       `json:"post_category_id" form:"post_category_id" binding:"omitempty"`
	Title          *string     `json:"title" form:"title" binding:"omitempty"`
	Date           *string     `json:"date" form:"date" binding:"omitempty"`
	UUIDs          []uuid.UUID `json:"uuids" form:"uuids" binding:"omitempty" ignore:"true"`
	Query          *string     `json:"query" form:"query" binding:"omitempty" ignore:"true"`
	models.Filter
}

func (p *CategoryFilter) ValidateForm(ctx *gin.Context) (*CategoryFilter, *errutil.Errors) {
	err := p.Filter.ValidateForm(ctx, p)
	return p, err
}

func (p *CategoryFilter) ValidateQuery(ctx *gin.Context) (*CategoryFilter, *errutil.Errors) {
	err := p.Filter.ValidateQuery(ctx, p)
	return p, err
}

func (p *CategoryFilter) PrintTemplateID() uint {
	if p.TemplateID == nil {
		return 0
	}
	return *p.TemplateID
}

func (p *CategoryFilter) PrintUserID() uint {
	if p.UserID == nil {
		return 0
	}
	return *p.UserID
}

func (p *CategoryFilter) PrintPostCategoryID() uint {
	if p.PostCategoryID == nil {
		return 0
	}
	return *p.PostCategoryID
}

func (p *CategoryFilter) GetURL() string {
	return string("post/categories?" + p.GetQueryString())
}

func (p *CategoryFilter) PrintID() string {
	if p.ID == nil {
		return ""
	}
	return strconv.Itoa(int(*p.ID))
}
