package provider

import (
	"github.com/axlle-com/blog/pkg/app/logger"
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	app "github.com/axlle-com/blog/pkg/app/service"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/repository"
	"github.com/axlle-com/blog/pkg/info_block/service"
)

type InfoBlockProvider interface {
	GetForResource(contracts.Resource) []contracts.InfoBlock
	GetAll() []contracts.InfoBlock
	SaveFromForm(g any, resource contracts.Resource) (contracts.InfoBlock, error)
	DeleteForResource(contracts.Resource) error
}

func NewProvider(
	infoBlockRepo repository.InfoBlockRepository,
	service *service.InfoBlockService,
) InfoBlockProvider {
	return &provider{
		infoBlockRepo: infoBlockRepo,
		service:       service,
	}
}

type provider struct {
	infoBlockRepo repository.InfoBlockRepository
	service       *service.InfoBlockService
}

func (p *provider) GetForResource(resource contracts.Resource) []contracts.InfoBlock {
	infoBlocks, err := p.infoBlockRepo.GetForResource(resource)
	collection := make([]contracts.InfoBlock, 0, len(infoBlocks))
	if err == nil {
		for _, infoBlock := range infoBlocks {
			collection = append(collection, infoBlock)
		}
		return collection
	}
	logger.Error(err)
	return nil
}

func (p *provider) DeleteForResource(resource contracts.Resource) (err error) {
	err = p.service.DeleteForResource(resource)
	if err != nil {
		return err
	}

	return nil
}

func (p *provider) GetAll() []contracts.InfoBlock {
	var collection []contracts.InfoBlock
	infoBlocks, err := p.infoBlockRepo.GetAll()
	if err == nil {
		for _, infoBlock := range infoBlocks {
			collection = append(collection, infoBlock)
		}
		return collection
	}
	logger.Error(err)
	return nil
}

func (p *provider) SaveFromForm(g any, resource contracts.Resource) (infoBlock contracts.InfoBlock, err error) {
	ib := app.LoadStruct(&models.InfoBlock{}, g).(*models.InfoBlock)
	if ib.ID == 0 {
		infoBlock, err = p.service.CreateInfoBlock(ib)
	} else {
		infoBlock, err = p.service.UpdateInfoBlock(ib)
	}

	if err != nil {
		return nil, err
	}

	err = p.service.Attach(resource, infoBlock)
	if err != nil {
		return nil, err
	}

	return infoBlock, nil
}
