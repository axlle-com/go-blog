package models

import (
	errorsForm "github.com/axlle-com/blog/app/errors"
	"github.com/gin-gonic/gin"
)

func NewBlockRequest() *BlockRequest {
	return &BlockRequest{}
}

type BlockRequest struct {
	ID          string            `gorm:"primaryKey" json:"id" form:"id" binding:"-"`
	TemplateID  string            `gorm:"index" json:"template_id" form:"template_id" binding:"omitempty"`
	Media       string            `gorm:"size:255" json:"media" form:"media" binding:"omitempty,max=255"`
	Title       string            `gorm:"size:255;not null" json:"title" form:"title" binding:"required,max=255"`
	Description string            `gorm:"type:text" json:"description" form:"description" binding:"omitempty"`
	Image       string            `gorm:"size:255" json:"image" form:"image" binding:"omitempty,max=255"`
	Galleries   []*GalleryRequest `json:"galleries" form:"galleries" binding:"omitempty"`
}

func (p *BlockRequest) ValidateForm(ctx *gin.Context) (*BlockRequest, *errorsForm.Errors) {
	err := ctx.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		return nil, &errorsForm.Errors{Message: "Форма не валидная!"}
	}
	if len(ctx.Request.PostForm) == 0 {
		return nil, &errorsForm.Errors{Message: "Форма не валидная!"}
	}
	if err := ctx.ShouldBind(&p); err != nil {
		errBind := errorsForm.ParseBindErrorToMap(err)
		return nil, errBind
	}
	return p, nil
}

func (p *BlockRequest) ValidateJSON(ctx *gin.Context) (*BlockRequest, *errorsForm.Errors) {
	if err := ctx.ShouldBindJSON(&p); err != nil {
		errBind := errorsForm.ParseBindErrorToMap(err)
		return nil, errBind
	}

	return p, nil
}
