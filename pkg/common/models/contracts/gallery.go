package contracts

import (
	"time"
)

type Gallery interface {
	GetID() uint
	GetTitle() *string
	GetDescription() *string
	GetSort() int
	GetImage() *string
	GetURL() *string
	GetDate() *time.Time
	GetImages() []GalleryImage
	Attach(Resource) error
}
