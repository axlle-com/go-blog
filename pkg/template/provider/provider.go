package provider

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/template/repository"
)

type TemplateProvider interface {
	GetAll() []contracts.Template
	GetAllIds() []uint
	GetByIDs(ids []uint) ([]contracts.Template, error)
	GetMapByIDs(ids []uint) (map[uint]contracts.Template, error)
}

func NewProvider(
	template repository.TemplateRepository,
) TemplateProvider {
	return &provider{
		templateRepo: template,
	}
}

type provider struct {
	templateRepo repository.TemplateRepository
}

func (p *provider) GetAll() []contracts.Template {
	all, err := p.templateRepo.GetAll()
	if err == nil {
		collection := make([]contracts.Template, 0, len(all))
		for _, t := range all {
			collection = append(collection, t)
		}
		return collection
	}
	logger.Error(err)
	return nil
}

func (p *provider) GetAllIds() []uint {
	t, err := p.templateRepo.GetAllIds()
	if err == nil {
		return t
	}
	logger.Error(err)
	return nil
}

func (p *provider) GetByIDs(ids []uint) ([]contracts.Template, error) {
	all, err := p.templateRepo.GetByIDs(ids)
	if err == nil {
		collection := make([]contracts.Template, 0, len(all))
		for _, t := range all {
			collection = append(collection, t)
		}
		return collection, nil
	}
	logger.Error(err)
	return nil, err
}

func (p *provider) GetMapByIDs(ids []uint) (map[uint]contracts.Template, error) {
	all, err := p.templateRepo.GetByIDs(ids)
	if err == nil {
		collection := make(map[uint]contracts.Template, len(all))
		for _, template := range all {
			collection[template.ID] = template
		}
		return collection, nil
	}
	logger.Error(err)
	return nil, err
}
