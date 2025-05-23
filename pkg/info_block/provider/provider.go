package provider

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	app "github.com/axlle-com/blog/app/service"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/service"
	"sync"
)

type InfoBlockProvider interface {
	GetForResource(contracts.Resource) []contracts.InfoBlock
	GetAll() []contracts.InfoBlock
	Attach(id uint, resource contracts.Resource) (infoBlocks []contracts.InfoBlock, err error)
	SaveForm(block any, resource contracts.Resource) (contracts.InfoBlock, error)
	SaveFormBatch(blocks []any, resource contracts.Resource) (infoBlock []contracts.InfoBlock, err error)
	DetachResource(contracts.Resource) error
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
	infoBlocks := p.blockService.GetForResource(resource)
	if infoBlocks == nil {
		return nil
	}

	collection := make([]contracts.InfoBlock, 0, len(infoBlocks))
	for _, infoBlock := range infoBlocks {
		collection = append(collection, infoBlock)
	}
	return collection
}

func (p *provider) DetachResource(resource contracts.Resource) (err error) {
	err = p.blockService.DetachResource(resource)
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

func (p *provider) Attach(id uint, resource contracts.Resource) (infoBlocks []contracts.InfoBlock, err error) {
	infoBlock, err := p.blockService.GetByID(id)
	if err != nil {
		return nil, err
	}

	err = p.blockService.Attach(resource, infoBlock)
	if err != nil {
		return nil, err
	}

	infoBlocks = p.GetForResource(resource)
	return infoBlocks, nil
}

func (p *provider) SaveForm(block any, resource contracts.Resource) (infoBlock contracts.InfoBlock, err error) {
	ib := app.LoadStruct(&models.InfoBlockResponse{}, block).(*models.InfoBlockResponse)

	infoBlock, err = p.blockService.GetByID(ib.GetID())
	if err != nil {
		return nil, err
	}
	ib.FromInterface(infoBlock)
	err = p.blockService.Attach(resource, ib)
	if err != nil {
		return nil, err
	}
	p.collectionService.AggregatesResponses([]*models.InfoBlockResponse{ib})
	return ib, nil
}

func (p *provider) SaveFormBatch(blocks []any, resource contracts.Resource) (infoBlock []contracts.InfoBlock, err error) {
	var wg sync.WaitGroup

	for _, block := range blocks {
		wg.Add(1)
		// Передаём block как параметр, чтобы избежать проблем замыкания
		go func(b any) {
			defer wg.Done()
			iBlock := app.LoadStruct(&models.InfoBlockResponse{}, b).(*models.InfoBlockResponse)
			if err := p.blockService.Attach(resource, iBlock); err != nil {
				logger.Error(err)
			}
		}(block)
	}

	wg.Wait()

	infoBlocks := p.GetForResource(resource)
	return infoBlocks, nil
}
