package models

import (
	"time"

	"github.com/axlle-com/blog/app/models/contract"
	"github.com/google/uuid"
)

type MenuItem struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	PublisherUUID *uuid.UUID `gorm:"type:uuid;index,using:hash"` // Publisher материал который доступен по URL
	MenuID        uint       `gorm:"index" json:"menu_id"`
	MenuItemID    *uint      `gorm:"index" json:"menu_item_id,omitempty"`
	PathLtree     string     `gorm:"type:ltree;column:path_ltree;not null" json:"-"`
	Title         string     `gorm:"size:100" json:"title"`
	URL           string     `gorm:"size:1000" json:"url"`
	Ico           *string    `gorm:"size:255" json:"ico,omitempty"`
	Sort          int        `gorm:"default:0" json:"sort,omitempty"`
	CreatedAt     *time.Time `gorm:"index" json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	DeletedAt     *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	Children  []*MenuItem
	Parent    *MenuItem
	Publisher contract.Publisher `gorm:"-" json:"publisher,omitempty"`
}

func (mi *MenuItem) GetTable() string {
	return "menu_items"
}

func (mi *MenuItem) Fields() []string {
	return []string{
		"UserID",
		"PublisherUUID",
		"MenuID",
		"MenuItemID",
		"PathLtree",
		"Title",
		"URL",
		"Ico",
		"Sort",
	}
}

func (mi *MenuItem) Creating() {
	mi.Saving()
}

func (mi *MenuItem) Updating() {
	mi.Saving()
}

func (mi *MenuItem) Deleting() bool {
	return true
}

func (mi *MenuItem) Saving() {}

func (mi *MenuItem) AdminAjaxFilterURL() string {
	return "/admin/ajax/menus/menus-items"
}
