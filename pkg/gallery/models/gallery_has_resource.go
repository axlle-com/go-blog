package models

import "github.com/google/uuid"

type GalleryHasResource struct {
	GalleryID    uint      `gorm:"index;not null"`
	ResourceUUID uuid.UUID `gorm:"type:uuid;index,using:hash" json:"uuid" form:"uuid" binding:"-"`
	Sort         int       `gorm:"index;not null"`
}

func (g *GalleryHasResource) GetTable() string {
	return "gallery_has_resources"
}
