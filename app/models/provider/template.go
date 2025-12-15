package provider

import "github.com/axlle-com/blog/app/models/contract"

type TemplateProvider interface {
	GetAll() []contract.Template
	GetByID(id uint) (contract.Template, error)
	GetAllIds() []uint
	GetByIDs(ids []uint) ([]contract.Template, error)
	GetMapByIDs(ids []uint) (map[uint]contract.Template, error)
	GetForResources(resource contract.Resource) ([]contract.Template, error)
	GetByNameAndResource(name string, resourceName string) (contract.Template, error)
}
