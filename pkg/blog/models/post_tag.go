package models

import (
	"fmt"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/google/uuid"
	"time"
)

type PostTag struct {
	ID              uint       `gorm:"primaryKey" json:"id" form:"id" binding:"-"`
	UUID            uuid.UUID  `gorm:"type:uuid;index,using:hash" json:"uuid" form:"uuid" binding:"-"`
	TemplateID      *uint      `gorm:"index" json:"template_id" form:"template_id" binding:"omitempty"`
	Name            string     `gorm:"size:10;not null;unique" json:"name" form:"name" binding:"required,max=10"`
	Title           *string    `gorm:"size:255" json:"title" form:"title" binding:"required,max=255"`
	Description     *string    `gorm:"type:text" json:"description" form:"description" binding:"omitempty"`
	Image           *string    `gorm:"size:255" json:"image" form:"image" binding:"omitempty,max=255"`
	MetaTitle       *string    `gorm:"size:100" json:"meta_title" form:"meta_title" binding:"omitempty,max=100"`
	MetaDescription *string    `gorm:"size:200" json:"meta_description" form:"meta_description" binding:"omitempty,max=200"`
	Alias           string     `gorm:"size:255;unique" json:"alias" form:"alias" binding:"omitempty,max=255"`
	URL             string     `gorm:"size:1000;unique" json:"url" form:"url" binding:"omitempty,max=1000"`
	CreatedAt       *time.Time `gorm:"index" json:"created_at,omitempty" form:"created_at" binding:"-" ignore:"true"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty" form:"updated_at" binding:"-" ignore:"true"`
	DeletedAt       *time.Time `gorm:"index" json:"deleted_at" form:"deleted_at" binding:"-" ignore:"true"`

	Galleries  []contracts.Gallery   `gorm:"-" json:"galleries" form:"galleries" binding:"-" ignore:"true"`
	InfoBlocks []contracts.InfoBlock `gorm:"-" json:"info_blocks" form:"info_blocks" binding:"-" ignore:"true"`
	Template   contracts.Template    `gorm:"-" json:"template" form:"template" binding:"-" ignore:"true"`
}

func (pt *PostTag) GetTable() string {
	return "post_tags"
}

func (pt *PostTag) GetID() uint {
	return pt.ID
}

func (pt *PostTag) GetUUID() uuid.UUID {
	return pt.UUID
}

func (pt *PostTag) GetName() string {
	return pt.GetTable()
}

func (pt *PostTag) GetURL() string {
	return pt.URL
}

func (pt *PostTag) GetTitle() string {
	if pt.Title == nil {
		return ""
	}
	return *pt.Title
}

func (pt *PostTag) AdminURL() string {
	if pt.ID == 0 {
		return "/admin/post-tags"
	}
	return fmt.Sprintf("/admin/post-tags/%d", pt.ID)
}

func (pt *PostTag) GetTemplateID() uint {
	var templateID uint
	if pt.TemplateID != nil {
		templateID = *pt.TemplateID
	}
	return templateID
}

func (pt *PostTag) GetTemplateName() string {
	if pt.Template != nil {
		if pt.Template.GetName() == "" {
			return fmt.Sprintf("%s.default", pt.GetTable())
		}
		return fmt.Sprintf("%s.%s", pt.GetTable(), pt.Template.GetName())
	}
	return fmt.Sprintf("%s.default", pt.GetTable())
}

func (pt *PostTag) GetGalleries() []contracts.Gallery {
	return pt.Galleries
}

func (pt *PostTag) GetDescription() string {
	return *pt.Description
}

func (pt *PostTag) GetImage() string {
	if pt.Image != nil {
		return *pt.Image
	}
	return ""
}

func (pt *PostTag) GetTemplateTitle() string {
	var title string
	if pt.Template != nil {
		title = pt.Template.GetTitle()
	}
	return title
}

func (pt *PostTag) GetAlias() string {
	return pt.Alias
}

func (pt *PostTag) Date() string {
	if pt.CreatedAt == nil {
		return ""
	}
	return pt.CreatedAt.Format("02.01.2006 15:04:05")
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
	pt.SetUUID()
}

func (pt *PostTag) SetUUID() {
	if pt.UUID == uuid.Nil {
		pt.UUID = uuid.New()
	}
}
