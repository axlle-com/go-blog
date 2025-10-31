package provider

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/template/models"
	"github.com/axlle-com/blog/pkg/template/repository"
)

type TemplateProvider interface {
	GetAll() []contract.Template
	GetByID(id uint) (contract.Template, error)
	GetAllIds() []uint
	GetByIDs(ids []uint) ([]contract.Template, error)
	GetMapByIDs(ids []uint) (map[uint]contract.Template, error)
	GetForResources(resource contract.Resource) ([]contract.Template, error)
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

func (p *provider) GetAll() []contract.Template {
	all, err := p.templateRepo.GetAll()
	if err == nil {
		collection := make([]contract.Template, 0, len(all))
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

func (p *provider) GetByID(id uint) (contract.Template, error) {
	model, err := p.templateRepo.GetByID(id)
	if err == nil {
		return model, nil
	}
	logger.Error(err)
	return nil, err
}

func (p *provider) GetByIDs(ids []uint) ([]contract.Template, error) {
	all, err := p.templateRepo.GetByIDs(ids)
	if err == nil {
		collection := make([]contract.Template, 0, len(all))
		for _, t := range all {
			collection = append(collection, t)
		}
		return collection, nil
	}
	logger.Error(err)
	return nil, err
}

func (p *provider) GetMapByIDs(ids []uint) (map[uint]contract.Template, error) {
	all, err := p.templateRepo.GetByIDs(ids)
	if err == nil {
		collection := make(map[uint]contract.Template, len(all))
		for _, template := range all {
			collection[template.ID] = template
		}
		return collection, nil
	}
	logger.Error(err)
	return nil, err
}

func (p *provider) GetForResources(resource contract.Resource) ([]contract.Template, error) {
	all, err := p.templateRepo.Filter(models.NewTemplateFilter().SetResourceName(resource.GetName()))
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	collection := make([]contract.Template, 0, len(all))
	for _, t := range all {
		collection = append(collection, t)
	}
	return collection, nil
}
