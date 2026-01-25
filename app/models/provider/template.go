package provider

import "github.com/axlle-com/blog/app/models/contract"

type TemplateProvider interface {
	GetAll() []contract.Template
	GetByID(id uint) (contract.Template, error)
	GetAllIds() []uint
	GetByIDs(ids []uint) ([]contract.Template, error)
	GetMapByIDs(ids []uint) (map[uint]contract.Template, error)
	GetMapByNames(names []string) (map[string]contract.Template, error)
	GetForResources(resource contract.Resource) ([]contract.Template, error)
	GetByName(name string) (contract.Template, error)
}
