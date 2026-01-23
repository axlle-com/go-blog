package request

import (
	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/pkg/blog/models"
	"github.com/gin-gonic/gin"
)

func NewCategoryRequest() *CategoryRequest {
	return &CategoryRequest{}
}

type CategoryRequest struct {
	ID                 uint    `json:"id" form:"id" binding:"omitempty"`
	UUID               string  `json:"uuid" form:"uuid" binding:"omitempty"`
	UserID             *uint   `json:"user_id" form:"user_id" binding:"omitempty"`
	TemplateID         *uint   `json:"template_id" form:"template_id" binding:"omitempty"`
	PostCategoryID     *uint   `json:"post_category_id" form:"post_category_id" binding:"omitempty"`
	MetaTitle          *string `json:"meta_title" form:"meta_title" binding:"omitempty,max=100"`
	MetaDescription    *string `json:"meta_description" form:"meta_description" binding:"omitempty,max=200"`
	Alias              string  `json:"alias" form:"alias" binding:"omitempty,max=255"`
	URL                string  `json:"url" form:"url" binding:"omitempty,max=1000"`
	IsPublished        *bool   `json:"is_published" form:"is_published" binding:"omitempty"`
	IsFavourites       *bool   `json:"is_favourites" form:"is_favourites" binding:"omitempty"`
	HasComments        *bool   `json:"has_comments" form:"has_comments" binding:"omitempty"`
	ShowImage          *bool   `json:"show_image" form:"show_image" binding:"omitempty"`
	InSitemap          *bool   `json:"in_sitemap" form:"in_sitemap" binding:"omitempty"`
	Title              string  `json:"title" form:"title" binding:"required,max=255"`
	TitleShort         *string `json:"title_short" form:"title_short" binding:"omitempty,max=155"`
	DescriptionPreview *string `json:"description_preview" form:"description_preview" binding:"omitempty"`
	Description        *string `json:"description" form:"description" binding:"omitempty"`
	ShowDate           string  `json:"show_date" form:"show_date" binding:"omitempty"`
	DatePub            string  `json:"date_pub,omitempty" time_format:"02.01.2006" form:"date_pub" binding:"omitempty"`
	DateEnd            string  `json:"date_end,omitempty" time_format:"02.01.2006" form:"date_end" binding:"omitempty"`
	Image              *string `json:"image" form:"image" binding:"omitempty,max=255"`
	Sort               *uint   `json:"sort" form:"sort" binding:"omitempty"`

	Galleries  []*GalleryRequest   `json:"galleries" form:"galleries" binding:"omitempty"`
	InfoBlocks []*InfoBlockRequest `json:"info_blocks" form:"info_blocks" binding:"omitempty"`
}

func (r *CategoryRequest) ValidateJSON(ctx *gin.Context) (*CategoryRequest, *errutil.Errors) {
	if err := ctx.ShouldBindJSON(&r); err != nil {
		return nil, errutil.NewErrors(err)
	}

	return r, nil
}

func (r *CategoryRequest) PreloadFromModel(model *models.PostCategory) {
	if model == nil {
		return
	}

	r.ID = model.ID
	r.UUID = model.UUID.String()
	r.UserID = model.UserID
}
