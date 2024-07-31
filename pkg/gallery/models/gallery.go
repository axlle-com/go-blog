package models

import "time"

type Gallery struct {
	ID           uint            `gorm:"primary_key" json:"id"`
	Title        *string         `gorm:"size:255" json:"title"`
	Description  *string         `gorm:"type:text" json:"description"`
	Sort         int             `gorm:"default:0" json:"sort"`
	Image        *string         `gorm:"size:255;" json:"image"`
	URL          *string         `gorm:"size:255" json:"url"`
	CreatedAt    *time.Time      `json:"created_at,omitempty"`
	UpdatedAt    *time.Time      `json:"updated_at,omitempty"`
	DeletedAt    *time.Time      `gorm:"index" json:"deleted_at,omitempty"`
	GalleryImage []*GalleryImage `json:"images,omitempty"`
}
