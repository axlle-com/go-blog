package contract

import "github.com/google/uuid"

// Publisher материал который доступен по URL
type Publisher interface {
	GetID() uint
	GetUUID() uuid.UUID
	GetURL() string
	GetTitle() string
	GetTable() string
}
