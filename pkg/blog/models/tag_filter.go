package models

import (
	"strconv"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/gin-gonic/gin"

	"github.com/axlle-com/blog/app/models"
	"github.com/google/uuid"
)

func NewTagFilter() *TagFilter {
	return &TagFilter{}
}

type TagFilter struct {
	ID         *uint       `json:"id" form:"id" binding:"omitempty"`
	TemplateID *uint       `json:"template_id" form:"template_id" binding:"omitempty"`
	Name       *string     `json:"name" form:"name" binding:"omitempty"`
	Title      *string     `json:"title" form:"title" binding:"omitempty"`
	Date       *string     `json:"date" form:"date" binding:"omitempty"`
	UUIDs      []uuid.UUID `json:"uuids" form:"uuids" binding:"omitempty" ignore:"true"`
	Query      *string     `json:"query" form:"query" binding:"omitempty" ignore:"true"`
	models.Filter
}

func (p *TagFilter) ValidateForm(ctx *gin.Context) (*TagFilter, *errutil.Errors) {
	err := p.Filter.ValidateForm(ctx, p)
	return p, err
}

func (p *TagFilter) ValidateQuery(ctx *gin.Context) (*TagFilter, *errutil.Errors) {
	err := p.Filter.ValidateQuery(ctx, p)
	return p, err
}

func (p *TagFilter) PrintTemplateID() uint {
	if p.TemplateID == nil {
		return 0
	}
	return *p.TemplateID
}

func (p *TagFilter) GetURL() string {
	return string("post/tags?" + p.GetQueryString())
}

func (p *TagFilter) PrintID() string {
	if p.ID == nil {
		return ""
	}
	return strconv.Itoa(int(*p.ID))
}
