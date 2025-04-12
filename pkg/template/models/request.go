package models

import (
	errorsForm "github.com/axlle-com/blog/app/errors"
	"github.com/gin-gonic/gin"
)

func NewTemplateRequest() *TemplateRequest {
	return &TemplateRequest{}
}

type TemplateRequest struct {
	ID      string `gorm:"primaryKey" json:"id" form:"id" binding:"-"`
	Title   string `gorm:"size:255;not null" json:"title" form:"title" binding:"required,max=255"`
	Name    string `gorm:"type:text" json:"name" form:"name" binding:"required,max=25"`
	Tabular string `gorm:"size:255" json:"tabular" form:"tabular" binding:"required,max=255"`
	JS      string `gorm:"type:text" json:"js" form:"js" binding:"omitempty"`
	CSS     string `gorm:"type:text" json:"css" form:"css" binding:"omitempty"`
}

func (p *TemplateRequest) ValidateForm(ctx *gin.Context) (*TemplateRequest, *errorsForm.Errors) {
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

func (p *TemplateRequest) ValidateJSON(ctx *gin.Context) (*TemplateRequest, *errorsForm.Errors) {
	if err := ctx.ShouldBindJSON(&p); err != nil {
		errBind := errorsForm.ParseBindErrorToMap(err)
		return nil, errBind
	}

	return p, nil
}
