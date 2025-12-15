package provider

import "github.com/axlle-com/blog/app/models/contract"

type AnalyticProvider interface {
	SaveForm(raw any) (contract.Analytic, error)
	GetAll() []contract.Analytic
	GetByID(id uint) (contract.Analytic, error)
	GetAllIds() []uint
	GetByIDs(ids []uint) ([]contract.Analytic, error)
	GetMapByIDs(ids []uint) (map[uint]contract.Analytic, error)
}
