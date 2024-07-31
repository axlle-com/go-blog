package models

import (
	"mime/multipart"
	"time"
)

type GalleryImage struct {
	ID           uint                  `gorm:"primary_key" json:"id"`
	GalleryID    uint                  `gorm:"not null;index" json:"gallery_id"`
	OriginalName string                `gorm:"size:255;not null" json:"original_name"`
	File         string                `gorm:"size:255;not null;unique" json:"file"`
	Title        *string               `gorm:"size:255" json:"title"`
	Description  *string               `gorm:"type:text" json:"description"`
	Sort         int                   `gorm:"default:0" json:"sort"`
	CreatedAt    *time.Time            `json:"created_at,omitempty"`
	UpdatedAt    *time.Time            `json:"updated_at,omitempty"`
	DeletedAt    *time.Time            `gorm:"index" json:"deleted_at,omitempty"`
	FileHeader   *multipart.FileHeader `gorm:"-"`
	Gallery      *Gallery
}
