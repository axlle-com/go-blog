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
	Gallery      *Gallery              `gorm:"-"`
}

func (*Image) TableName() string {
	return "gallery_images"
}

func (i *Image) GetID() uint {
	return i.ID
}

func (i *Image) GetGalleryID() uint {
	return i.GalleryID
}

func (i *Image) GetTitle() *string {
	return i.Title
}

func (i *Image) GetDescription() *string {
	return i.Description
}

func (i *Image) GetSort() int {
	return i.Sort
}

func (i *Image) GetFile() string {
	return i.File
}

func (i *Image) GetDate() *time.Time {
	return i.CreatedAt
}

func (i *Image) GetGallery() contracts.Gallery {
	return i.Gallery
}

func (i *Image) Deleted() error {
	err := file.DeleteFile(i.File)
	if err != nil {
		return err
	}
	count := ImageRepo().CountForGallery(i.GalleryID)
	if count == 0 {
		err = GalleryRepo().Delete(i.GalleryID)
		if err != nil {
			return err
		}
	}
	return nil
}
