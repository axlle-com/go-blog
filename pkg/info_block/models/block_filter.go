package models

import (
	"strconv"

	"github.com/axlle-com/blog/app/models"
	"github.com/google/uuid"
)

func NewInfoBlockFilter() *InfoBlockFilter {
	return &InfoBlockFilter{}
}

type InfoBlockFilter struct {
	ID             *uint       `json:"id" form:"id" binding:"omitempty"`
	UUIDs          []uuid.UUID `json:"uuid" form:"uuid" binding:"-"`
	TemplateID     *uint       `json:"template_id" form:"template_id" binding:"omitempty"`
	UserID         *uint       `json:"user_id" form:"user_id" binding:"omitempty"`
	PostCategoryID *uint       `json:"post_category_id" form:"post_category_id" binding:"omitempty"`
	Title          *string     `json:"title" form:"title" binding:"omitempty"`
	Date           *string     `json:"date" form:"date" binding:"omitempty"`
	ResourceUUID   *uuid.UUID  `json:"resource_uuid" form:"resource_uuid" binding:"omitempty"`
	RelationID     *uint       `json:"relation_id" form:"relation_id" binding:"omitempty"`
	RelationIDs    []uint      `json:"relation_ids" form:"relation_ids" binding:"omitempty"`
	InfoBlockIDs   []uint      `json:"info_block_ids" form:"info_block_ids" binding:"omitempty"`
	models.Filter
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
