package contracts

import (
	"github.com/google/uuid"
)

type InfoBlock interface {
	GetID() uint
	GetUUID() uuid.UUID
	GetTemplateID() uint
	GetTemplateTitle() string
	GetTemplateName() string
	GetTitle() string
	GetDescription() string
	GetImage() string
	GetMedia() string
	GetGalleries() []Gallery
	GetPosition() string
	GetPositions() []string
	GetSort() int
	GetRelationID() uint
}
