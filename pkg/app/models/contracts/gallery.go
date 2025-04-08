package contracts

import (
	"github.com/google/uuid"
	"time"
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
