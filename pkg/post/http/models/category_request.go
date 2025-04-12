package models

import (
	errorsForm "github.com/axlle-com/blog/app/errors"
	"github.com/gin-gonic/gin"
)

func NewCategoryRequest() *CategoryRequest {
	return &CategoryRequest{}
}

type CategoryRequest struct {
	ID                 string `json:"id" form:"id" binding:"omitempty"`
	UUID               string `json:"uuid" form:"uuid" binding:"omitempty"`
	UserID             string `json:"user_id" form:"user_id" binding:"omitempty"`
	TemplateID         string `json:"template_id" form:"template_id" binding:"omitempty"`
	PostCategoryID     string `json:"post_category_id" form:"post_category_id" binding:"omitempty"`
	MetaTitle          string `json:"meta_title" form:"meta_title" binding:"omitempty,max=100"`
	MetaDescription    string `json:"meta_description" form:"meta_description" binding:"omitempty,max=200"`
	Alias              string `json:"alias" form:"alias" binding:"omitempty,max=255"`
	URL                string `json:"url" form:"url" binding:"omitempty,max=1000"`
	IsPublished        string `json:"is_published" form:"is_published" binding:"omitempty"`
	IsFavourites       string `json:"is_favourites" form:"is_favourites" binding:"omitempty"`
	HasComments        string `json:"has_comments" form:"has_comments" binding:"omitempty"`
	ShowImagePost      string `json:"show_image_post" form:"show_image_post" binding:"omitempty"`
	ShowImageCategory  string `json:"show_image_category" form:"show_image_category" binding:"omitempty"`
	InSitemap          string `json:"in_sitemap" form:"in_sitemap" binding:"omitempty"`
	Media              string `json:"media" form:"media" binding:"omitempty,max=255"`
	Title              string `json:"title" form:"title" binding:"required,max=255"`
	TitleShort         string `json:"title_short" form:"title_short" binding:"omitempty,max=155"`
	DescriptionPreview string `json:"description_preview" form:"description_preview" binding:"omitempty"`
	Description        string `json:"description" form:"description" binding:"omitempty"`
	ShowDate           string `json:"show_date" form:"show_date" binding:"omitempty"`
	DatePub            string `json:"date_pub,omitempty" time_format:"02.01.2006" form:"date_pub" binding:"omitempty"`
	DateEnd            string `json:"date_end,omitempty" time_format:"02.01.2006" form:"date_end" binding:"omitempty"`
	Image              string `json:"image" form:"image" binding:"omitempty,max=255"`
	Sort               string `json:"sort" form:"sort" binding:"omitempty"`

	Galleries  []*GalleryRequest   `json:"galleries" form:"galleries" binding:"omitempty"`
	InfoBlocks []*InfoBlockRequest `json:"info_blocks" form:"info_blocks" binding:"omitempty"`
}

func (p *CategoryRequest) ValidateForm(ctx *gin.Context) (*CategoryRequest, *errorsForm.Errors) {
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

func (p *CategoryRequest) ValidateJSON(ctx *gin.Context) (*CategoryRequest, *errorsForm.Errors) {
	if err := ctx.ShouldBindJSON(&p); err != nil {
		errBind := errorsForm.ParseBindErrorToMap(err)
		return nil, errBind
	}

	return p, nil
}
