package contracts

import "time"

type GalleryImage interface {
	GetID() uint
	GetGalleryID() uint
	GetTitle() *string
	GetDescription() *string
	GetSort() int
	GetFile() string
	GetFilePath() string
	GetDate() *time.Time
	GetGallery() Gallery
}
