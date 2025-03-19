package models

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type MenuItem struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	ResourceUUID uuid.UUID  `gorm:"type:uuid;index,using:hash"`
	MenuID       *uint      `gorm:"index" json:"menu_id,required"`
	MenuItemID   *uint      `gorm:"index" json:"menu_item_id,omitempty"`
	MatPath      string     `gorm:"size:1000;index" json:"-"`
	Title        *string    `gorm:"size:100" json:"title,omitempty"`
	URL          string     `gorm:"size:1000;unique" json:"url"`
	IsPublished  *bool      `gorm:"default:true" json:"is_published,omitempty"`
	Ico          *string    `gorm:"size:255" json:"ico,omitempty"`
	Sort         *uint      `gorm:"default:0" json:"sort,omitempty"`
	CreatedAt    *time.Time `gorm:"index" json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (c *MenuItem) GetMenuID() uint {
	var menuID uint
	if c.MenuID != nil {
		menuID = *c.MenuID
	}
	return menuID
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

func (c *MenuItem) GetMatPath() string {
	return fmt.Sprintf("%s.%d", c.MatPath, c.ID)
}
