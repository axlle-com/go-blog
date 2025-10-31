package provider

import (
	"sync"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	apppPovider "github.com/axlle-com/blog/app/models/provider"
	app "github.com/axlle-com/blog/app/service"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/service"
	"github.com/google/uuid"
)

func NewProvider(
	blockService *service.InfoBlockService,
	collectionService *service.InfoBlockCollectionService,
) apppPovider.InfoBlockProvider {
	return &provider{
		blockService:      blockService,
		collectionService: collectionService,
	}
}

type provider struct {
	blockService      *service.InfoBlockService
	collectionService *service.InfoBlockCollectionService
}

func (p *provider) GetForResourceUUID(resourceUUID string) []contract.InfoBlock {
	newUUID, err := uuid.Parse(resourceUUID)
	if err != nil {
		logger.Errorf("[info_block][provider] invalid resource_uuid: %v", err)
		return nil
	}

	filter := models.NewInfoBlockFilter()
	filter.ResourceUUID = &newUUID
	infoBlocks := p.blockService.GetForResourceByFilter(filter)
	if infoBlocks == nil {
		return nil
	}

	collection := make([]contract.InfoBlock, 0, len(infoBlocks))
	for _, infoBlock := range infoBlocks {
		collection = append(collection, infoBlock)
	}
	return collection
}

func (p *provider) DetachResourceUUID(resourceUUID string) (err error) {
	newUUID, err := uuid.Parse(resourceUUID)
	if err != nil {
		logger.Errorf("[info_block][provider] invalid resource_uuid: %v", err)
		return nil
	}

	err = p.blockService.DeleteByResourceUUID(newUUID)
	if err != nil {
		return err
	}

	return nil
}

func (p *provider) GetAll() []contract.InfoBlock {
	var collection []contract.InfoBlock
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

func (p *provider) Attach(id uint, resourceUUID string) (infoBlocks []contract.InfoBlock, err error) {
	newUUID, err := uuid.Parse(resourceUUID)
	if err != nil {
		logger.Errorf("[info_block][provider] invalid resource_uuid: %v", err)
		return
	}

	infoBlock, err := p.blockService.FindByID(id)
	if err != nil {
		return nil, err
	}

	err = p.blockService.Attach(newUUID, infoBlock)
	if err != nil {
		return nil, err
	}

	infoBlocks = p.GetForResourceUUID(resourceUUID)
	return infoBlocks, nil
}

func (p *provider) SaveForm(block any, resourceUUID string) (infoBlock contract.InfoBlock, err error) {
	newUUID, err := uuid.Parse(resourceUUID)
	if err != nil {
		logger.Errorf("[info_block][provider] invalid resource_uuid: %v", err)
		return
	}

	ib := app.LoadStruct(&models.InfoBlockResponse{}, block).(*models.InfoBlockResponse)

	infoBlock, err = p.blockService.FindByID(ib.GetID())
	if err != nil {
		return nil, err
	}
	ib.FromInterface(infoBlock)
	err = p.blockService.Attach(newUUID, ib)
	if err != nil {
		return nil, err
	}
	p.collectionService.AggregatesResponses([]*models.InfoBlockResponse{ib})
	return ib, nil
}

func (p *provider) SaveFormBatch(blocks []any, resourceUUID string) (infoBlock []contract.InfoBlock, err error) {
	newUUID, err := uuid.Parse(resourceUUID)
	if err != nil {
		logger.Errorf("[info_block][provider] invalid resource_uuid: %v", err)
		return
	}

	var wg sync.WaitGroup

	for _, block := range blocks {
		wg.Add(1)
		// Передаём block как параметр, чтобы избежать проблем замыкания
		go func(b any) {
			defer wg.Done()
			iBlock := app.LoadStruct(&models.InfoBlockResponse{}, b).(*models.InfoBlockResponse)
			if err := p.blockService.Attach(newUUID, iBlock); err != nil {
				logger.Error(err)
			}
		}(block)
	}

	wg.Wait()

	infoBlocks := p.GetForResourceUUID(resourceUUID)
	return infoBlocks, nil
}
