package models

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"mime/multipart"
	"os"
	"path/filepath"
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

func (gi *GalleryImage) GetFilePath() string {
	absPath, err := filepath.Abs("src" + gi.File)
	if err != nil {
		logger.New().Error(err)
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		logger.New().Error(err)
	}
	return absPath
}

func (gi *GalleryImage) GetDate() *time.Time {
	return gi.CreatedAt
}

func (gi *GalleryImage) GetGallery() contracts.Gallery {
	return gi.Gallery
}

func (gi *GalleryImage) Deleted() error {
	err := os.Remove(gi.GetFilePath())
	if err != nil {
		return err
	}
	count := NewGalleryImageRepository().CountForGallery(gi.GalleryID)
	if count == 0 {
		err = NewGalleryRepository().Delete(gi.GalleryID)
		if err != nil {
			return err
		}
	}
	return nil
}
