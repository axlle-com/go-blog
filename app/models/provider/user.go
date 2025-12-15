package provider

import (
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/google/uuid"
)

type UserProvider interface {
	GetAll() []contract.User
	GetAllIds() []uint
	GetByID(id uint) (contract.User, error)
	GetByUUID(uuid uuid.UUID) (contract.User, error)
	GetByIDs(ids []uint) ([]contract.User, error)
	GetMapByIDs(ids []uint) (map[uint]contract.User, error)
	GetMapByUUIDs(uuids []uuid.UUID) (map[uuid.UUID]contract.User, error)
	Create(contract.User) (contract.User, error)
}
