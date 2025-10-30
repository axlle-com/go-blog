package contract

import (
	"time"

	"github.com/google/uuid"
)

type Gallery interface {
	GetID() uint
	GetResourceUUID() uuid.UUID
	GetTitle() *string
	GetDescription() *string
	GetSort() int
	GetPosition() string
	GetImage() *string
	GetURL() *string
	GetDate() *time.Time
	GetImages() []Image
}
