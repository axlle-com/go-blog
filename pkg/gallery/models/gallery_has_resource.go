package models

type GalleryHasResource struct {
	GalleryID  uint   `gorm:"index;not null"`
	Resource   string `gorm:"index;not null;size:255"`
	ResourceID uint   `gorm:"index;not null"`
}
