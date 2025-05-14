package models

import (
	"github.com/axlle-com/blog/app/errutil"
	. "github.com/axlle-com/blog/app/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

func NewInfoBlockFilter() *InfoBlockFilter {
	return &InfoBlockFilter{}
}

type InfoBlockFilter struct {
	ID             *uint   `json:"id" form:"id" binding:"omitempty"`
	TemplateID     *uint   `json:"template_id" form:"template_id" binding:"omitempty"`
	UserID         *uint   `json:"user_id" form:"user_id" binding:"omitempty"`
	PostCategoryID *uint   `json:"post_category_id" form:"post_category_id" binding:"omitempty"`
	Title          *string `json:"title" form:"title" binding:"omitempty"`
	Date           *string `json:"date" form:"date" binding:"omitempty"`
	Filter
}

func (p *InfoBlockFilter) ValidateForm(ctx *gin.Context) (*InfoBlockFilter, *errutil.Errors) {
	err := p.Filter.ValidateForm(ctx, p)
	return p, err
}

func (p *InfoBlockFilter) ValidateQuery(ctx *gin.Context) (*InfoBlockFilter, *errutil.Errors) {
	err := p.Filter.ValidateQuery(ctx, p)
	return p, err
}

func (p *InfoBlockFilter) PrintTemplateID() uint {
	if p.TemplateID == nil {
		return 0
	}
	return *p.TemplateID
}

func (p *InfoBlockFilter) PrintUserID() uint {
	if p.UserID == nil {
		return 0
	}
	return *p.UserID
}

func (p *InfoBlockFilter) PrintPostCategoryID() uint {
	if p.PostCategoryID == nil {
		return 0
	}
	return *p.PostCategoryID
}

func (p *InfoBlockFilter) GetURL() string {
	return string("info-blocks?" + p.GetQueryString())
}

func (p *InfoBlockFilter) PrintID() string {
	if p.ID == nil {
		return ""
	}
	return strconv.Itoa(int(*p.ID))
}
