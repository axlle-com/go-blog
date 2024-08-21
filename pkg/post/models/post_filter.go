package models

import (
	. "github.com/axlle-com/blog/pkg/common/models"
	"github.com/gin-gonic/gin"
	"net/http"
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

func (p *PostFilter) ValidateForm(ctx *gin.Context) *PostFilter {
	err := p.Filter.ValidateForm(ctx, p)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errors":  err.Errors,
			"message": err.Message,
		})
		ctx.Abort()
		return nil
	}
	return p
}

func (p *PostFilter) ValidateQuery(ctx *gin.Context) *PostFilter {
	err := p.Filter.ValidateQuery(ctx, p)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errors":  err.Errors,
			"message": err.Message,
		})
		ctx.Abort()
		return nil
	}
	return p
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
	return string("posts?" + p.GetQueryString())
}
