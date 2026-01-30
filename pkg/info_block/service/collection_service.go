package service

import (
	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/repository"
)

type CollectionService struct {
	api                        *api.Api
	infoBlockRepo              repository.InfoBlockRepository
	resourceRepo               repository.InfoBlockHasResourceRepository
	collectionAggregateService *CollectionAggregateService
}

func NewCollectionService(
	api *api.Api,
	infoBlockRepo repository.InfoBlockRepository,
	resourceRepo repository.InfoBlockHasResourceRepository,
	collectionAggregateService *CollectionAggregateService,
) *CollectionService {
	return &CollectionService{
		api:                        api,
		infoBlockRepo:              infoBlockRepo,
		resourceRepo:               resourceRepo,
		collectionAggregateService: collectionAggregateService,
	}
}

func (s *CollectionService) GetAll() ([]*models.InfoBlock, error) {
	return s.infoBlockRepo.GetAll()
}

func (s *CollectionService) GetRoots() ([]*models.InfoBlock, error) {
	return s.infoBlockRepo.GetRoots()
}

func (s *CollectionService) GetForResourceByFilter(filter *models.InfoBlockFilter) []*models.InfoBlockResponse {
	infoBlocks, err := s.infoBlockRepo.GetForResourceByFilter(filter)
	if err != nil {
		logger.Errorf("[info_block][CollectionService][GetForResourceByFilter] Error: %v", err)
		return nil
	}

	return s.collectionAggregateService.AggregatesResponses(infoBlocks)
}

func (s *CollectionService) GetAllForParent(parent *models.InfoBlock) ([]*models.InfoBlock, error) {
	return s.infoBlockRepo.GetAllForParent(parent)
}

func (s *CollectionService) WithPaginate(paginator contract.Paginator, filter *models.InfoBlockFilter) ([]*models.InfoBlock, error) {
	return s.infoBlockRepo.WithPaginate(paginator, filter)
}

func (s *CollectionService) Aggregates(infoBlocks []*models.InfoBlock) []*models.InfoBlock {
	return s.collectionAggregateService.Aggregates(infoBlocks)
}
