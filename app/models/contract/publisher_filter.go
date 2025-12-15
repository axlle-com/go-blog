package contract

import "github.com/google/uuid"

type PublisherFilter interface {
	GetUUIDs() []uuid.UUID
	GetQuery() string
	GetURL() *string
}
