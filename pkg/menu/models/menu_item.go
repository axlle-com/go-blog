package models

import (
	"time"

	"github.com/google/uuid"
)

type MenuItem struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	PublisherUUID uuid.UUID  `gorm:"type:uuid;index,using:hash"` // Publisher материал который доступен по URL
	MenuID        uint       `gorm:"index" json:"menu_id"`
	MenuItemID    *uint      `gorm:"index" json:"menu_item_id,omitempty"`
	Path          string     `gorm:"size:1000;index" json:"-"`
	Title         string     `gorm:"size:100" json:"title"`
	URL           string     `gorm:"size:1000;unique" json:"url"`
	IsPublished   bool       `gorm:"default:true" json:"is_published,omitempty"`
	Ico           *string    `gorm:"size:255" json:"ico,omitempty"`
	Sort          int        `gorm:"default:0" json:"sort,omitempty"`
	CreatedAt     *time.Time `gorm:"index" json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	DeletedAt     *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (c *MenuItem) Creating() {
	c.Saving()
}

func (c *MenuItem) Updating() {
	c.Saving()
}

func (c *MenuItem) Deleting() bool {
	return true
}

func (c *MenuItem) Saving() {
}
