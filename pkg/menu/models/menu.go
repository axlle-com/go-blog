package models

import (
	"time"
)

type Menu struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	TemplateID  *uint      `gorm:"index" json:"template_id"`
	Name        *string    `gorm:"size:100" json:"name,omitempty"`
	IsPublished *bool      `gorm:"default:true" json:"is_published,omitempty"`
	Ico         *string    `gorm:"size:255" json:"ico,omitempty"`
	Sort        *uint      `gorm:"default:0" json:"sort,omitempty"`
	CreatedAt   *time.Time `gorm:"index" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (c *Menu) GetTemplateID() uint {
	var templateID uint
	if c.TemplateID != nil {
		templateID = *c.TemplateID
	}
	return templateID
}

func (c *Menu) Creating() {
	c.Saving()
}

func (c *Menu) Updating() {
	c.Saving()
}

func (c *Menu) Deleting() bool {
	return true
}

func (c *Menu) Saving() {
}
