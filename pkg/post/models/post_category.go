package models

import (
	"github.com/google/uuid"
	"time"
)

type PostCategory struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	UUID               uuid.UUID  `gorm:"type:uuid;index,using:hash" json:"uuid" form:"uuid" binding:"-"`
	TemplateID         *uint      `gorm:"index" json:"template_id,omitempty"`
	PostCategoryID     *uint      `gorm:"index" json:"post_category_id,omitempty"`
	MetaTitle          *string    `gorm:"size:100" json:"meta_title,omitempty"`
	MetaDescription    *string    `gorm:"size:200" json:"meta_description,omitempty"`
	Alias              string     `gorm:"size:255;unique" json:"alias"`
	URL                string     `gorm:"size:1000;unique" json:"url"`
	IsPublished        *bool      `gorm:"default:true" json:"is_published,omitempty"`
	IsFavourites       *bool      `gorm:"default:false" json:"is_favourites,omitempty"`
	InSitemap          *bool      `gorm:"default:true" json:"in_sitemap,omitempty"`
	Image              *string    `gorm:"size:255" json:"image,omitempty"`
	ShowImage          *bool      `gorm:"default:true" json:"show_image,omitempty"`
	Title              string     `gorm:"size:255;not null" json:"title"`
	TitleShort         *string    `gorm:"size:150" json:"title_short,omitempty"`
	Description        *string    `gorm:"type:text" json:"description,omitempty"`
	DescriptionPreview *string    `gorm:"type:text" json:"description_preview,omitempty"`
	Sort               *uint      `gorm:"default:0" json:"sort,omitempty"`
	CreatedAt          *time.Time `json:"created_at,omitempty"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
	DeletedAt          *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (c *PostCategory) GetUUID() uuid.UUID {
	return c.UUID
}

func (c *PostCategory) SetUUID() {
	if c.UUID == uuid.Nil {
		c.UUID = uuid.New()
	}
}
