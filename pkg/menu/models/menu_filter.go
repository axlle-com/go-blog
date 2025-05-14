package models

import (
	"github.com/axlle-com/blog/app/errutil"
	. "github.com/axlle-com/blog/app/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

func NewMenuFilter() *MenuFilter {
	return &MenuFilter{}
}

type MenuFilter struct {
	ID             *uint   `json:"id" form:"id" binding:"omitempty"`
	TemplateID     *uint   `json:"template_id" form:"template_id" binding:"omitempty"`
	UserID         *uint   `json:"user_id" form:"user_id" binding:"omitempty"`
	MenuCategoryID *uint   `json:"post_category_id" form:"post_category_id" binding:"omitempty"`
	Title          *string `json:"title" form:"title" binding:"omitempty"`
	Date           *string `json:"date" form:"date" binding:"omitempty"`
	Filter
}

func (p *MenuFilter) ValidateForm(ctx *gin.Context) (*MenuFilter, *errutil.Errors) {
	err := p.Filter.ValidateForm(ctx, p)
	return p, err
}

func (p *MenuFilter) ValidateQuery(ctx *gin.Context) (*MenuFilter, *errutil.Errors) {
	err := p.Filter.ValidateQuery(ctx, p)
	return p, err
}

func (p *MenuFilter) PrintTemplateID() uint {
	if p.TemplateID == nil {
		return 0
	}
	return *p.TemplateID
}

func (p *MenuFilter) PrintUserID() uint {
	if p.UserID == nil {
		return 0
	}
	return *p.UserID
}

func (p *MenuFilter) PrintMenuCategoryID() uint {
	if p.MenuCategoryID == nil {
		return 0
	}
	return *p.MenuCategoryID
}

func (p *MenuFilter) GetURL() string {
	if p.GetQueryString() == "" {
		return "posts"
	}
	return string("posts?" + p.GetQueryString())
}

func (p *MenuFilter) PrintID() string {
	if p.ID == nil {
		return ""
	}
	return strconv.Itoa(int(*p.ID))
}
