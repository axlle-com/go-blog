package models

import (
	"github.com/axlle-com/blog/app/errors"
	. "github.com/axlle-com/blog/app/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

func NeTemplateFilter() *TemplateFilter {
	return &TemplateFilter{}
}

type TemplateFilter struct {
	ID      *uint   `json:"id" form:"id" binding:"omitempty"`
	UserID  *uint   `json:"user_id" form:"user_id" binding:"omitempty"`
	Title   *string `json:"title" form:"title" binding:"omitempty"`
	Name    *string `json:"name" form:"name" binding:"omitempty"`
	Tabular *string `json:"tabular" form:"tabular" binding:"omitempty"`
	Date    *string `json:"date" form:"date" binding:"omitempty"`
	Filter
}

func (p *TemplateFilter) ValidateForm(ctx *gin.Context) (*TemplateFilter, *errors.Errors) {
	err := p.Filter.ValidateForm(ctx, p)
	return p, err
}

func (p *TemplateFilter) ValidateQuery(ctx *gin.Context) (*TemplateFilter, *errors.Errors) {
	err := p.Filter.ValidateQuery(ctx, p)
	return p, err
}

func (p *TemplateFilter) PrintUserID() uint {
	if p.UserID == nil {
		return 0
	}
	return *p.UserID
}

func (p *TemplateFilter) GetURL() string {
	return string("templates?" + p.GetQueryString())
}

func (p *TemplateFilter) PrintID() string {
	if p.ID == nil {
		return ""
	}
	return strconv.Itoa(int(*p.ID))
}
