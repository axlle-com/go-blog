package provider

import (
	"github.com/axlle-com/blog/app/logger"
	contracts2 "github.com/axlle-com/blog/app/models/contracts"
	app "github.com/axlle-com/blog/app/service"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/service"
)

type InfoBlockProvider interface {
	GetForResource(contracts2.Resource) []contracts2.InfoBlock
	GetAll() []contracts2.InfoBlock
	Attach(id uint, resource contracts2.Resource) (infoBlocks []contracts2.InfoBlock, err error)
	SaveForm(block any, resource contracts2.Resource) (contracts2.InfoBlock, error)
	SaveFormBatch(blocks []any, resource contracts2.Resource) (infoBlock []contracts2.InfoBlock, err error)
	DeleteForResource(contracts2.Resource) error
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

func (p *provider) GetForResource(resource contracts2.Resource) []contracts2.InfoBlock {
	infoBlocks := p.blockService.GetForResource(resource)
	if infoBlocks == nil {
		return nil
	}

	collection := make([]contracts2.InfoBlock, 0, len(infoBlocks))
	for _, infoBlock := range infoBlocks {
		collection = append(collection, infoBlock)
	}
	return collection
}

func (p *provider) DeleteForResource(resource contracts2.Resource) (err error) {
	err = p.blockService.DeleteForResource(resource)
	if err != nil {
		return err
	}

	return nil
}

func (p *provider) GetAll() []contracts2.InfoBlock {
	var collection []contracts2.InfoBlock
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

func (p *provider) Attach(id uint, resource contracts2.Resource) (infoBlocks []contracts2.InfoBlock, err error) {
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

func (p *provider) SaveForm(block any, resource contracts2.Resource) (infoBlock contracts2.InfoBlock, err error) {
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

func (p *provider) SaveFormBatch(blocks []any, resource contracts2.Resource) (infoBlock []contracts2.InfoBlock, err error) {
	var blocksIDs []uint
	infoBlockResponses := make([]*models.InfoBlockResponse, 0, len(blocks))
	for _, block := range blocks {
		iBlock := app.LoadStruct(&models.InfoBlockResponse{}, block).(*models.InfoBlockResponse)
		infoBlockResponses = append(infoBlockResponses, iBlock)
		blocksIDs = append(blocksIDs, iBlock.GetID())
	}

	for _, iBlock := range infoBlockResponses {
		err = p.blockService.Attach(resource, iBlock)
		if err != nil {
			return nil, err
		}
	}

	infoBlocks := p.GetForResource(resource)
	return infoBlocks, nil
}
