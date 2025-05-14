package contracts

import "github.com/google/uuid"

type PostTag interface {
	GetID() uint
	GetUUID() uuid.UUID
	GetTemplateID() uint
	GetTitle() string
	GetDescription() string
	GetImage() string
	GetGalleries() []Gallery
}
