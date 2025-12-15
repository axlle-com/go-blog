package provider

import "github.com/axlle-com/blog/app/models/contract"

type MessageProvider interface {
	GetAll() []contract.Message
	GetByID(id uint) (contract.Message, error)
	GetAllIds() []uint
	GetByIDs(ids []uint) ([]contract.Message, error)
	GetMapByIDs(ids []uint) (map[uint]contract.Message, error)
}
