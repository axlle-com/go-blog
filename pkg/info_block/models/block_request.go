package models

import (
	"github.com/axlle-com/blog/app/errutil"
	"github.com/gin-gonic/gin"
)

func NewBlockRequest() *BlockRequest {
	return &BlockRequest{}
}

type BlockRequest struct {
	ID          uint              `gorm:"primaryKey" json:"id" form:"id" binding:"-"`
	TemplateID  *uint             `gorm:"index" json:"template_id" form:"template_id" binding:"omitempty"`
	InfoBlockID *uint             `gorm:"index" json:"info_block_id" form:"info_block_id" binding:"omitempty"`
	Media       *string           `gorm:"size:255" json:"media" form:"media" binding:"omitempty,max=255"`
	Title       string            `gorm:"size:255;not null" json:"title" form:"title" binding:"required,max=255"`
	Description *string           `gorm:"type:text" json:"description" form:"description" binding:"omitempty"`
	Image       *string           `gorm:"size:255" json:"image" form:"image" binding:"omitempty,max=255"`
	Galleries   []*GalleryRequest `json:"galleries" form:"galleries" binding:"omitempty"`
}

func (p *BlockRequest) ValidateJSON(ctx *gin.Context) (*BlockRequest, *errutil.Errors) {
	if err := ctx.ShouldBindJSON(&p); err != nil {
		return nil, errutil.NewErrors(err)
	}

	return p, nil
}
