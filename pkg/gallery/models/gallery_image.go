package models

import (
	"github.com/axlle-com/blog/pkg/common/models/contracts"
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

func (gi *GalleryImage) GetID() uint {
	return gi.ID
}

func (gi *GalleryImage) GetGalleryID() uint {
	return gi.GalleryID
}

func (gi *GalleryImage) GetTitle() *string {
	return gi.Title
}

func (gi *GalleryImage) GetDescription() *string {
	return gi.Description
}

func (gi *GalleryImage) GetSort() int {
	return gi.Sort
}

func (gi *GalleryImage) GetFile() string {
	return gi.File
}

func (gi *GalleryImage) GetDate() *time.Time {
	return gi.CreatedAt
}

func (gi *GalleryImage) GetGallery() contracts.Gallery {
	return gi.Gallery
}

func (gi *GalleryImage) AfterDelete() {

}
