package provider

import (
	"github.com/axlle-com/blog/pkg/app/logger"
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	"github.com/axlle-com/blog/pkg/template/repository"
)

type TemplateProvider interface {
	GetAll() []contracts.Template
	GetAllIds() []uint
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
	ts, err := p.templateRepo.GetAll()
	if err == nil {
		var collection []contracts.Template
		for _, t := range ts {
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
