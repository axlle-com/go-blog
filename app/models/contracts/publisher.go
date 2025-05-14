package contracts

import "github.com/google/uuid"

type Publisher interface {
	GetID() uint
	GetUUID() uuid.UUID
	GetURL() string
	GetTitle() string
	GetTable() string
}
