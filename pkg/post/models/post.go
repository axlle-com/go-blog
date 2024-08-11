package models

import (
	"github.com/axlle-com/blog/pkg/gallery/models"
	"time"
)

type Post struct {
	ID                 uint              `gorm:"primaryKey" json:"id" form:"id" binding:"-"`
	UserID             *uint             `gorm:"index" json:"user_id" form:"user_id" binding:"omitempty"`
	TemplateID         *uint             `gorm:"index" json:"template_id" form:"template_id" binding:"omitempty"`
	PostCategoryID     *uint             `gorm:"index" json:"post_category_id" form:"post_category_id" binding:"omitempty"`
	MetaTitle          *string           `gorm:"size:100" json:"meta_title" form:"meta_title" binding:"omitempty,max=100"`
	MetaDescription    *string           `gorm:"size:200" json:"meta_description" form:"meta_description" binding:"omitempty,max=200"`
	Alias              string            `gorm:"size:255;unique" json:"alias" form:"alias" binding:"required,max=255"`
	URL                string            `gorm:"size:1000;unique" json:"url" form:"url" binding:"required,max=1000"`
	IsPublished        bool              `gorm:"not null;default:true" json:"is_published" form:"is_published" binding:"omitempty"`
	IsFavourites       bool              `gorm:"not null;default:false" json:"is_favourites" form:"is_favourites" binding:"omitempty"`
	HasComments        bool              `gorm:"not null;default:false" json:"has_comments" form:"has_comments" binding:"omitempty"`
	ShowImagePost      bool              `gorm:"not null;default:true" json:"show_image_post" form:"show_image_post" binding:"omitempty"`
	ShowImageCategory  bool              `gorm:"not null;default:true" json:"show_image_category" form:"show_image_category" binding:"omitempty"`
	MakeWatermark      bool              `gorm:"not null;default:false" json:"make_watermark" form:"make_watermark" binding:"omitempty"`
	InSitemap          bool              `gorm:"not null;default:true" json:"in_sitemap" form:"in_sitemap" binding:"omitempty"`
	Media              *string           `gorm:"size:255" json:"media" form:"media" binding:"omitempty,max=255"`
	Title              string            `gorm:"size:255;not null" json:"title" form:"title" binding:"required,max=255"`
	TitleShort         *string           `gorm:"size:155" json:"title_short" form:"title_short" binding:"omitempty,max=155"`
	DescriptionPreview *string           `gorm:"type:text" json:"description_preview" form:"description_preview" binding:"omitempty"`
	Description        *string           `gorm:"type:text" json:"description" form:"description" binding:"omitempty"`
	ShowDate           bool              `gorm:"not null;default:true" json:"show_date" form:"show_date" binding:"omitempty"`
	DatePub            *time.Time        `json:"date_pub,omitempty" time_format:"02.01.2006" form:"date_pub" binding:"omitempty"`
	DateEnd            *time.Time        `json:"date_end,omitempty" time_format:"02.01.2006" form:"date_end" binding:"omitempty"`
	Image              *string           `gorm:"size:255" json:"image" form:"image" binding:"omitempty,max=255"`
	Hits               uint              `gorm:"not null;default:0" json:"hits" form:"hits" binding:"-"`
	Sort               int               `gorm:"not null;default:0" json:"sort" form:"sort" binding:"omitempty"`
	Stars              float32           `gorm:"not null;default:0.0" json:"stars" form:"stars" binding:"-"`
	CreatedAt          *time.Time        `json:"created_at,omitempty" form:"created_at" binding:"-"`
	UpdatedAt          *time.Time        `json:"updated_at,omitempty" form:"updated_at" binding:"-"`
	DeletedAt          *time.Time        `gorm:"index" json:"deleted_at" form:"deleted_at" binding:"-"`
	Galleries          []*models.Gallery `gorm:"-" json:"galleries" form:"galleries" binding:"-"`
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

func (p *Post) GetID() uint {
	return p.ID
}

func (p *Post) GetResource() string {
	return "posts"
}
