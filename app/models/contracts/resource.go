package contracts

import "github.com/google/uuid"

type Resource interface {
	GetUUID() uuid.UUID
	GetName() string
	GetTemplateName() string
}
