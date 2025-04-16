package contracts

import "github.com/google/uuid"

type User interface {
	GetID() uint
	GetUUID() uuid.UUID
	GetFirstName() string
	GetLastName() string
	GetPatronymic() string
	GetPhone() string
	GetEmail() string
	GetStatus() int8
	GetRoles() []string
	GetPermissions() []string
}
