package models

import (
	errorsForm "github.com/axlle-com/blog/pkg/common/errors"
	. "github.com/axlle-com/blog/pkg/common/models"
	. "github.com/axlle-com/blog/pkg/post/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func NewPostRequest() *PostRequest {
	return &PostRequest{}
}

type PostRequest struct {
	TemplateID         *uint      `json:"template_id" form:"template_id" binding:"omitempty"`
	PostCategoryID     *uint      `json:"post_category_id" form:"post_category_id" binding:"omitempty"`
	MetaTitle          *string    `json:"meta_title" form:"meta_title" binding:"omitempty,max=100"`
	MetaDescription    *string    `json:"meta_description" form:"meta_description" binding:"omitempty,max=200"`
	Alias              *string    `json:"alias" form:"alias" binding:"omitempty,max=255"`
	IsPublished        bool       `json:"is_published" form:"is_published" binding:"omitempty"`
	IsFavourites       bool       `json:"is_favourites" form:"is_favourites" binding:"omitempty"`
	HasComments        bool       `json:"has_comments" form:"has_comments" binding:"omitempty"`
	ShowImagePost      bool       `json:"show_image_post" form:"show_image_post" binding:"omitempty"`
	ShowImageCategory  bool       `json:"show_image_category" form:"show_image_category" binding:"omitempty"`
	MakeWatermark      bool       `json:"make_watermark" form:"make_watermark" binding:"omitempty"`
	InSitemap          bool       `json:"in_sitemap" form:"in_sitemap" binding:"omitempty"`
	Media              *string    `json:"media" form:"media" binding:"omitempty,max=255"`
	Title              string     `json:"title" form:"title" binding:"required,max=255"`
	TitleShort         *string    `json:"title_short" form:"title_short" binding:"omitempty,max=155"`
	DescriptionPreview *string    `json:"description_preview" form:"description_preview" binding:"omitempty"`
	Description        *string    `json:"description" form:"description" binding:"omitempty"`
	ShowDate           bool       `json:"show_date" form:"show_date" binding:"omitempty"`
	DatePub            *time.Time `json:"date_pub,omitempty" time_format:"02.01.2006" form:"date_pub" binding:"omitempty"`
	DateEnd            *time.Time `json:"date_end,omitempty" time_format:"02.01.2006" form:"date_end" binding:"omitempty"`
	Sort               int        `json:"sort" form:"sort" binding:"omitempty"`
}

func (p *PostRequest) ValidateForm(ctx *gin.Context) *Post {
	err := ctx.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Форма не валидная"})
		ctx.Abort()
		return nil
	}
	var form *Post
	if err := ctx.ShouldBind(&form); err != nil {
		errBind := errorsForm.ParseBindErrorToMap(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errors":  errBind,
			"message": "Ошибки валидации",
		})
		ctx.Abort()
		return nil
	}
	SetEmptyStringPointersToNil(form)
	return form
}
