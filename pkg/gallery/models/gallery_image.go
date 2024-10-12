package models

import (
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"github.com/axlle-com/blog/pkg/file"
	"mime/multipart"
	"time"
)

type Image struct {
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

func (*Image) TableName() string {
	return "gallery_images"
}

func (gi *Image) GetID() uint {
	return gi.ID
}

func (gi *Image) GetGalleryID() uint {
	return gi.GalleryID
}

func (gi *Image) GetTitle() *string {
	return gi.Title
}

func (gi *Image) GetDescription() *string {
	return gi.Description
}

func (gi *Image) GetSort() int {
	return gi.Sort
}

func (gi *Image) GetFile() string {
	return gi.File
}

func (gi *Image) GetDate() *time.Time {
	return gi.CreatedAt
}

func (gi *Image) GetGallery() contracts.Gallery {
	return gi.Gallery
}

func (gi *Image) Deleted() error {
	err := file.DeleteFile(gi.File)
	if err != nil {
		return err
	}
	count := ImageRepo().CountForGallery(gi.GalleryID)
	if count == 0 {
		err = GalleryRepo().Delete(gi.GalleryID)
		if err != nil {
			return err
		}
	}
	return nil
}
