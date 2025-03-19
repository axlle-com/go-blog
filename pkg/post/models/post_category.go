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
	LeftSet            int        `gorm:"index" json:"-"`
	RightSet           int        `gorm:"index" json:"-"`
	Level              int        `gorm:"index" json:"-"`
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
	CreatedAt          *time.Time `gorm:"index" json:"created_at,omitempty"`
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

func (c *PostCategory) SetAlias() {
	if c.Alias != "" {
		return
	}
	if c.UUID == uuid.Nil {
		c.SetUUID()
	}
	c.Alias = c.UUID.String()
}

func (c *PostCategory) GetCategoryID() uint {
	var categoryID uint
	if c.PostCategoryID != nil {
		categoryID = *c.PostCategoryID
	}
	return categoryID
}

func (c *PostCategory) GetTemplateID() uint {
	var templateID uint
	if c.TemplateID != nil {
		templateID = *c.TemplateID
	}
	return templateID
}

func (c *PostCategory) GetID() uint {
	return c.ID
}

func (c *PostCategory) GetTable() string {
	return "post_categories"
}

func (c *PostCategory) Creating() {
	c.Saving()
}

func (c *PostCategory) Updating() {
	c.Saving()
}

func (c *PostCategory) Deleting() bool {
	return true
}

func (c *PostCategory) Saving() {
	c.SetUUID()
	c.SetAlias()
	c.setTitleShort()
	c.setURL()
}

func (c *PostCategory) setURL() {
	if c.Alias != "" {
		c.URL = "/" + c.Alias
	}
}

func (c *PostCategory) setTitleShort() {
	if c.TitleShort == nil {
		return
	}
	if *c.TitleShort == "" {
		c.TitleShort = nil
	}
}
