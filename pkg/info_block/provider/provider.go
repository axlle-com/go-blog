package provider

import (
	"github.com/axlle-com/blog/pkg/app/logger"
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	app "github.com/axlle-com/blog/pkg/app/service"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/service"
)

type InfoBlockProvider interface {
	GetForResource(contracts.Resource) []contracts.InfoBlock
	GetAll() []contracts.InfoBlock
	SaveFromForm(g any, resource contracts.Resource) (contracts.InfoBlock, error)
	DeleteForResource(contracts.Resource) error
}

func NewProvider(
	blockService *service.InfoBlockService,
	collectionService *service.InfoBlockCollectionService,
) InfoBlockProvider {
	return &provider{
		blockService:      blockService,
		collectionService: collectionService,
	}
}

type provider struct {
	blockService      *service.InfoBlockService
	collectionService *service.InfoBlockCollectionService
}

func (p *provider) GetForResource(resource contracts.Resource) []contracts.InfoBlock {
	infoBlocks, err := p.blockService.GetForResource(resource)
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
	err = p.blockService.DeleteForResource(resource)
	if err != nil {
		return err
	}

	return nil
}

func (p *provider) GetAll() []contracts.InfoBlock {
	var collection []contracts.InfoBlock
	infoBlocks, err := p.collectionService.GetAll()
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
		infoBlock, err = p.blockService.Create(ib, nil)
	} else {
		infoBlock, err = p.blockService.Update(ib)
	}

	if err != nil {
		return nil, err
	}

	err = p.blockService.Attach(resource, infoBlock)
	if err != nil {
		return nil, err
	}

	return infoBlock, nil
}
