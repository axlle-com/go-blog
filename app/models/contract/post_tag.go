package contract

import "github.com/google/uuid"

type PostTag interface {
	GetID() uint
	GetUUID() uuid.UUID
	GetTemplateName() string
	GetTitle() string
	GetDescription() string
	GetImage() string
	GetGalleries() []Gallery
}
