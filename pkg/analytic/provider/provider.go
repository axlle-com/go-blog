package provider

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	app "github.com/axlle-com/blog/app/service"
	"github.com/axlle-com/blog/pkg/analytic/models"
	"github.com/axlle-com/blog/pkg/analytic/service"
)

type AnalyticProvider interface {
	SaveForm(raw any) (contracts.Analytic, error)
	GetAll() []contracts.Analytic
	GetByID(id uint) (contracts.Analytic, error)
	GetAllIds() []uint
	GetByIDs(ids []uint) ([]contracts.Analytic, error)
	GetMapByIDs(ids []uint) (map[uint]contracts.Analytic, error)
}

func NewAnalyticProvider(
	analyticService *service.AnalyticService,
	analyticCollectionService *service.AnalyticCollectionService,
) AnalyticProvider {
	return &provider{
		analyticService:           analyticService,
		analyticCollectionService: analyticCollectionService,
	}
}

type provider struct {
	analyticService           *service.AnalyticService
	analyticCollectionService *service.AnalyticCollectionService
}

func (p *provider) SaveForm(raw any) (contracts.Analytic, error) {
	temp := app.LoadStruct(&models.Analytic{}, raw).(*models.Analytic)

	analytic, err := p.analyticService.Create(temp)
	if err != nil {
		return nil, err
	}
	return analytic, nil
}

func (p *provider) GetAll() []contracts.Analytic {
	all, err := p.analyticCollectionService.GetAll()
	if err == nil {
		collection := make([]contracts.Analytic, 0, len(all))
		for _, t := range all {
			collection = append(collection, t)
		}
		return collection
	}
	logger.Error(err)
	return nil
}

func (p *provider) GetAllIds() []uint {
	t, err := p.analyticCollectionService.GetAllIds()
	if err == nil {
		return t
	}
	logger.Error(err)
	return nil
}

func (p *provider) GetByID(id uint) (contracts.Analytic, error) {
	model, err := p.analyticService.GetByID(id)
	if err == nil {
		return model, nil
	}
	logger.Error(err)
	return nil, err
}

func (p *provider) GetByIDs(ids []uint) ([]contracts.Analytic, error) {
	all, err := p.analyticCollectionService.GetByIDs(ids)
	if err == nil {
		collection := make([]contracts.Analytic, 0, len(all))
		for _, t := range all {
			collection = append(collection, t)
		}
		return collection, nil
	}
	logger.Error(err)
	return nil, err
}

func (p *provider) GetMapByIDs(ids []uint) (map[uint]contracts.Analytic, error) {
	all, err := p.analyticCollectionService.GetByIDs(ids)
	if err == nil {
		collection := make(map[uint]contracts.Analytic, len(all))
		for _, analytic := range all {
			collection[analytic.ID] = analytic
		}
		return collection, nil
	}
	logger.Error(err)
	return nil, err
}
