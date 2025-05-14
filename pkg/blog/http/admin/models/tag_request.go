package models

import (
	errorsForm "github.com/axlle-com/blog/app/errutil"
	"github.com/gin-gonic/gin"
)

func NewTagRequest() *TagRequest {
	return &TagRequest{}
}

type TagRequest struct {
	ID              string `json:"id" form:"id" binding:"omitempty"`
	UUID            string `json:"uuid" form:"uuid" binding:"omitempty"`
	TemplateID      string `json:"template_id" form:"template_id" binding:"omitempty"`
	Name            string `json:"name" form:"name" binding:"required,max=10"`
	Title           string `json:"title" form:"title" binding:"required,max=255"`
	Description     string `json:"description" form:"description" binding:"omitempty"`
	Image           string `json:"image" form:"image" binding:"omitempty,max=255"`
	Sort            string `json:"sort" form:"sort" binding:"omitempty"`
	MetaTitle       string `json:"meta_title" form:"meta_title" binding:"omitempty,max=100"`
	MetaDescription string `json:"meta_description" form:"meta_description" binding:"omitempty,max=200"`
	Alias           string `json:"alias" form:"alias" binding:"omitempty,max=255"`
	URL             string `json:"url" form:"url" binding:"omitempty,max=1000"`

	Galleries  []*GalleryRequest   `json:"galleries" form:"galleries" binding:"omitempty"`
	InfoBlocks []*InfoBlockRequest `json:"info_blocks" form:"info_blocks" binding:"omitempty"`
}

func (p *TagRequest) ValidateForm(ctx *gin.Context) (*TagRequest, *errorsForm.Errors) {
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

func (p *TagRequest) ValidateJSON(ctx *gin.Context) (*TagRequest, *errorsForm.Errors) {
	if err := ctx.ShouldBindJSON(&p); err != nil {
		errBind := errorsForm.ParseBindErrorToMap(err)
		return nil, errBind
	}

	return p, nil
}
