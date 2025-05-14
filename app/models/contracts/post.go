package contracts

import (
	"github.com/google/uuid"
)

type Post interface {
	GetID() uint
	GetUUID() uuid.UUID
	GetTitle() string
	GetDescription() *string
	GetImage() *string
	GetURL() string
	Date() string
}
