package models

import (
	"github.com/axlle-com/blog/app/models/contracts"
	"time"
)

type PostTag struct {
	ID              uint       `gorm:"primaryKey" json:"id" form:"id" binding:"-"`
	TemplateID      *uint      `gorm:"index" json:"template_id" form:"template_id" binding:"omitempty"`
	Name            string     `gorm:"size:10;not null;unique" json:"name" form:"name" binding:"required,max=10"`
	Title           *string    `gorm:"size:255;not null" json:"title" form:"title" binding:"required,max=255"`
	Description     *string    `gorm:"type:text" json:"description" form:"description" binding:"omitempty"`
	Image           *string    `gorm:"size:255" json:"image" form:"image" binding:"omitempty,max=255"`
	MetaTitle       *string    `gorm:"size:100" json:"meta_title" form:"meta_title" binding:"omitempty,max=100"`
	MetaDescription *string    `gorm:"size:200" json:"meta_description" form:"meta_description" binding:"omitempty,max=200"`
	Alias           string     `gorm:"size:255;unique" json:"alias" form:"alias" binding:"omitempty,max=255"`
	URL             string     `gorm:"size:1000;unique" json:"url" form:"url" binding:"omitempty,max=1000"`
	CreatedAt       *time.Time `gorm:"index" json:"created_at,omitempty" form:"created_at" binding:"-" ignore:"true"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty" form:"updated_at" binding:"-" ignore:"true"`
	DeletedAt       *time.Time `gorm:"index" json:"deleted_at" form:"deleted_at" binding:"-" ignore:"true"`

	Galleries []contracts.Gallery `gorm:"-" json:"galleries" form:"galleries" binding:"-" ignore:"true"`
}

func (pt *PostTag) GetID() uint {
	return pt.ID
}

func (pt *PostTag) GetTemplateID() uint {
	var templateID uint
	if pt.TemplateID != nil {
		templateID = *pt.TemplateID
	}
	return templateID
}

func (pt *PostTag) GetTable() string {
	return "post_tags"
}

func (pt *PostTag) GetTitle() string {
	return *pt.Title
}

func (pt *PostTag) GetAlias() string {
	return pt.Alias
}

func (pt *PostTag) GetDescription() string {
	return *pt.Description
}

func (pt *PostTag) GetImage() string {
	return *pt.Image
}

func (pt *PostTag) GetGalleries() []contracts.Gallery {
	return pt.Galleries
}

func (pt *PostTag) setURL() {
	if pt.Alias != "" {
		pt.URL = "/" + pt.Alias
	}
}

func (pt *PostTag) Creating() {
	pt.Saving()
}

func (pt *PostTag) Updating() {
	pt.Saving()
}

func (pt *PostTag) Deleting() bool {
	return true
}

func (pt *PostTag) Saving() {
	pt.setURL()
}
