package models

import (
	"github.com/axlle-com/blog/app/errors"
	. "github.com/axlle-com/blog/app/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

func NewAnalyticFilter() *AnalyticFilter {
	return &AnalyticFilter{}
}

type AnalyticFilter struct {
	ID           *uint   `json:"id" form:"id" binding:"omitempty"`
	UserID       *uint   `json:"user_id" form:"user_id" binding:"omitempty"`
	Title        *string `json:"title" form:"title" binding:"omitempty"`
	Name         *string `json:"name" form:"name" binding:"omitempty"`
	ResourceName *string `json:"resource_name" form:"resource_name" binding:"omitempty"`
	Date         *string `json:"date" form:"date" binding:"omitempty"`
	Filter
}

func (p *AnalyticFilter) ValidateForm(ctx *gin.Context) (*AnalyticFilter, *errors.Errors) {
	err := p.Filter.ValidateForm(ctx, p)
	return p, err
}

func (p *AnalyticFilter) ValidateQuery(ctx *gin.Context) (*AnalyticFilter, *errors.Errors) {
	err := p.Filter.ValidateQuery(ctx, p)
	return p, err
}

func (p *AnalyticFilter) PrintUserID() uint {
	if p.UserID == nil {
		return 0
	}
	return *p.UserID
}

func (p *AnalyticFilter) GetURL() string {
	return string("templates?" + p.GetQueryString())
}

func (p *AnalyticFilter) PrintID() string {
	if p.ID == nil {
		return ""
	}
	return strconv.Itoa(int(*p.ID))
}
