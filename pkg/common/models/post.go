package models

import (
	"time"
)

type Post struct {
	ID                 uint          `gorm:"primaryKey" json:"id"`
	UserID             *uint         `gorm:"index" json:"user_id"`
	TemplateID         *uint         `gorm:"index" json:"template_id,omitempty"`
	PostCategoryID     *uint         `gorm:"index" json:"post_category_id,omitempty"`
	MetaTitle          *string       `gorm:"size:100" json:"meta_title,omitempty"`
	MetaDescription    *string       `gorm:"size:200" json:"meta_description,omitempty"`
	Alias              string        `gorm:"size:255;unique" json:"alias"`
	URL                string        `gorm:"size:1000;unique" json:"url"`
	IsPublished        bool          `gorm:"not null;default:true" json:"is_published,omitempty"`
	IsFavourites       bool          `gorm:"not null;default:false" json:"is_favourites,omitempty"`
	HasComments        bool          `gorm:"not null;default:false" json:"has_comments,omitempty"`
	ShowImagePost      bool          `gorm:"not null;default:true" json:"show_image_post,omitempty"`
	ShowImageCategory  bool          `gorm:"not null;default:true" json:"show_image_category,omitempty"`
	MakeWatermark      bool          `gorm:"not null;default:false" json:"make_watermark,omitempty"`
	InSitemap          bool          `gorm:"not null;default:true" json:"in_sitemap,omitempty"`
	Media              *string       `gorm:"size:255" json:"media,omitempty"`
	Title              string        `gorm:"size:255;not null" json:"title"`
	TitleShort         *string       `gorm:"size:155" json:"title_short,omitempty"`
	DescriptionPreview *string       `gorm:"type:text" json:"description_preview,omitempty"`
	Description        *string       `gorm:"type:text" json:"description,omitempty"`
	ShowDate           bool          `gorm:"not null;default:true" json:"show_date,omitempty"`
	DatePub            *time.Time    `json:"date_pub,omitempty"`
	DateEnd            *time.Time    `json:"date_end,omitempty"`
	Image              *string       `gorm:"size:255" json:"image,omitempty"`
	Hits               uint          `gorm:"not null;default:0" json:"hits,omitempty"`
	Sort               int           `gorm:"not null;default:0" json:"sort,omitempty"`
	Stars              float32       `gorm:"not null;default:0.0" json:"stars,omitempty"`
	CreatedAt          *time.Time    `json:"created_at,omitempty"`
	UpdatedAt          *time.Time    `json:"updated_at,omitempty"`
	DeletedAt          *time.Time    `gorm:"index" json:"deleted_at,omitempty"`
	Galleries          []interface{} `gorm:"-" json:"galleries,omitempty"`
}

func (p *Post) GetCategoryID() uint {
	var categoryID uint
	if p.PostCategoryID != nil {
		categoryID = *p.PostCategoryID
	}
	return categoryID
}

func (p *Post) GetTemplateID() uint {
	var templateID uint
	if p.TemplateID != nil {
		templateID = *p.TemplateID
	}
	return templateID
}
