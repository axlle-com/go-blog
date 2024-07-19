package models

import (
	"time"
)

type Post struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	TemplateID         *uint      `gorm:"index" json:"template_id,omitempty"`
	PostCategoryID     *uint      `gorm:"index" json:"post_category_id,omitempty"`
	MetaTitle          *string    `gorm:"size:100" json:"meta_title,omitempty"`
	MetaDescription    *string    `gorm:"size:200" json:"meta_description,omitempty"`
	Alias              string     `gorm:"size:255;unique" json:"alias"`
	URL                string     `gorm:"size:1000;unique" json:"url"`
	IsPublished        *bool      `gorm:"default:true" json:"is_published,omitempty"`
	IsFavourites       *bool      `gorm:"default:false" json:"is_favourites,omitempty"`
	IsComments         *bool      `gorm:"default:false" json:"is_comments,omitempty"`
	IsImagePost        *bool      `gorm:"default:true" json:"is_image_post,omitempty"`
	IsImageCategory    *bool      `gorm:"default:true" json:"is_image_category,omitempty"`
	IsWatermark        *bool      `gorm:"default:true" json:"is_watermark,omitempty"`
	IsSitemap          *bool      `gorm:"default:true" json:"is_sitemap,omitempty"`
	Media              *string    `gorm:"size:255" json:"media,omitempty"`
	Title              string     `gorm:"size:255;not null" json:"title"`
	TitleShort         *string    `gorm:"size:155" json:"title_short,omitempty"`
	PreviewDescription *string    `gorm:"type:text" json:"preview_description,omitempty"`
	Description        *string    `gorm:"type:text" json:"description,omitempty"`
	ShowDate           *bool      `gorm:"default:true" json:"show_date,omitempty"`
	DatePub            *time.Time `json:"date_pub,omitempty"`
	DateEnd            *time.Time `json:"date_end,omitempty"`
	ControlDatePub     *bool      `gorm:"default:false" json:"control_date_pub,omitempty"`
	ControlDateEnd     *bool      `gorm:"default:false" json:"control_date_end,omitempty"`
	Image              *string    `gorm:"size:255" json:"image,omitempty"`
	Hits               *uint      `gorm:"default:0" json:"hits,omitempty"`
	Sort               *int       `gorm:"default:0" json:"sort,omitempty"`
	Stars              *float32   `gorm:"default:0.0" json:"stars,omitempty"`
	CreatedAt          *time.Time `json:"created_at,omitempty"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
	DeletedAt          *time.Time `json:"deleted_at,;index;omitempty"`
}
