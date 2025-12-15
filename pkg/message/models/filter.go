package models

import (
	"strconv"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/models"
	"github.com/gin-gonic/gin"
)

func NewMessageFilter() *MessageFilter {
	return &MessageFilter{}
}

type MessageFilter struct {
	ID           *uint   `json:"id" form:"id" binding:"omitempty"`
	UserID       *uint   `json:"user_id" form:"user_id" binding:"omitempty"`
	Title        *string `json:"title" form:"title" binding:"omitempty"`
	Name         *string `json:"name" form:"name" binding:"omitempty"`
	ResourceName *string `json:"resource_name" form:"resource_name" binding:"omitempty"`
	Date         *string `json:"date" form:"date" binding:"omitempty"`
	models.Filter
}

func (p *MessageFilter) ValidateQuery(ctx *gin.Context) (*MessageFilter, *errutil.Errors) {
	err := p.Filter.ValidateQuery(ctx, p)
	return p, err
}

func (p *MessageFilter) PrintUserID() uint {
	if p.UserID == nil {
		return 0
	}
	return *p.UserID
}

func (p *MessageFilter) GetURL() string {
	return string("templates?" + p.GetQueryString())
}

func (p *MessageFilter) PrintID() string {
	if p.ID == nil {
		return ""
	}
	return strconv.Itoa(int(*p.ID))
}
