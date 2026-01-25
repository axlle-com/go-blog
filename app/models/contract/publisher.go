package contract

import "github.com/google/uuid"

const MaxShotTitle = 25

// Publisher материал который доступен по URL
type Publisher interface {
	GetID() uint
	GetUUID() uuid.UUID
	GetURL() string
	GetTitle() string
	GetImage() string
	GetMetaTitle() string
	GetMetaDescription() string
	GetTemplateName() string
	GetTable() string
}
