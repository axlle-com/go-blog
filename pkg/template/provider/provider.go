package provider

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"github.com/axlle-com/blog/pkg/template/repository"
)

type Template interface {
	GetAll() []contracts.Template
	GetAllIds() []uint
}

func Provider() Template {
	return &provider{}
}

type provider struct {
}

func (p *provider) GetAll() []contracts.Template {
	ts, err := repository.NewRepo().GetAll()
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
	t, err := repository.NewRepo().GetAllIds()
	if err == nil {
		return t
	}
	logger.Error(err)
	return nil
}
