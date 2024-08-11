package models

import "time"

type PostForm struct {
	UserID             *uint      `form:"user_id" binding:"required"`
	TemplateID         *uint      `form:"template_id" binding:"required"`
	PostCategoryID     *uint      `form:"post_category_id" binding:"required"`
	MetaTitle          *string    `form:"meta_title" binding:"omitempty,max=100"`
	MetaDescription    *string    `form:"meta_description" binding:"omitempty,max=200"`
	Alias              string     `form:"alias" binding:"required,max=255"`
	URL                string     `form:"url" binding:"required,max=1000"`
	IsPublished        bool       `form:"is_published" binding:"omitempty"`
	IsFavourites       bool       `form:"is_favourites" binding:"omitempty"`
	HasComments        bool       `form:"has_comments" binding:"omitempty"`
	ShowImagePost      bool       `form:"show_image_post" binding:"omitempty"`
	ShowImageCategory  bool       `form:"show_image_category" binding:"omitempty"`
	MakeWatermark      bool       `form:"make_watermark" binding:"omitempty"`
	InSitemap          bool       `form:"in_sitemap" binding:"omitempty"`
	Media              *string    `form:"media" binding:"omitempty,max=255"`
	Title              string     `form:"title" binding:"required,max=255"`
	TitleShort         *string    `form:"title_short" binding:"omitempty,max=155"`
	DescriptionPreview *string    `form:"description_preview" binding:"omitempty"`
	Description        *string    `form:"description" binding:"omitempty"`
	ShowDate           bool       `form:"show_date" binding:"omitempty"`
	DatePub            *time.Time `form:"date_pub" time_format:"02.01.2006" binding:"omitempty"`
	DateEnd            *time.Time `form:"date_end" time_format:"02.01.2006" binding:"omitempty"`
	Image              *string    `form:"image" binding:"required,max=255"`
	Hits               uint       `form:"hits" binding:"omitempty"`
	Sort               int        `form:"sort" binding:"omitempty"`
	Stars              float32    `form:"stars" binding:"omitempty"`
}
