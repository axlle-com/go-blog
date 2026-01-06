package service

import (
	"sync"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/repository"
	"github.com/google/uuid"
)

type InfoBlockCollectionService struct {
	infoBlockRepo repository.InfoBlockRepository
	resourceRepo  repository.InfoBlockHasResourceRepository
	api           *api.Api
}

func NewInfoBlockCollectionService(
	infoBlockRepo repository.InfoBlockRepository,
	resourceRepo repository.InfoBlockHasResourceRepository,
	api *api.Api,
) *InfoBlockCollectionService {
	return &InfoBlockCollectionService{
		infoBlockRepo: infoBlockRepo,
		resourceRepo:  resourceRepo,
		api:           api,
	}
}

func (s *InfoBlockCollectionService) GetAll() ([]*models.InfoBlock, error) {
	return s.infoBlockRepo.GetAll()
}

func (s *InfoBlockCollectionService) GetRoots() ([]*models.InfoBlock, error) {
	return s.infoBlockRepo.GetRoots()
}

func (s *InfoBlockCollectionService) GetForResourceByFilter(filter *models.InfoBlockFilter) []*models.InfoBlockResponse {
	infoBlocks, err := s.infoBlockRepo.GetForResourceByFilter(filter)
	if err != nil {
		logger.Error(err)
		return nil
	}

	return s.AggregatesResponses(infoBlocks)
}

func (s *InfoBlockCollectionService) GetAllForParent(parent *models.InfoBlock) ([]*models.InfoBlock, error) {
	return s.infoBlockRepo.GetAllForParent(parent)
}

func (s *InfoBlockCollectionService) WithPaginate(paginator contract.Paginator, filter *models.InfoBlockFilter) ([]*models.InfoBlock, error) {
	return s.infoBlockRepo.WithPaginate(paginator, filter)
}

func (s *InfoBlockCollectionService) Aggregates(infoBlocks []*models.InfoBlock) []*models.InfoBlock {
	var templateIDs []uint
	var userIDs []uint

	templateIDsMap := make(map[uint]bool)
	userIDsMap := make(map[uint]bool)

	for _, infoBlock := range infoBlocks {
		if infoBlock.TemplateID != nil {
			id := *infoBlock.TemplateID
			if !templateIDsMap[id] {
				templateIDs = append(templateIDs, id)
				templateIDsMap[id] = true
			}
		}
		if infoBlock.UserID != nil {
			id := *infoBlock.UserID
			if !userIDsMap[id] {
				userIDs = append(userIDs, id)
				userIDsMap[id] = true
			}
		}
	}

	var wg sync.WaitGroup

	var users map[uint]contract.User
	var templates map[uint]contract.Template

	wg.Add(2)

	go func() {
		defer wg.Done()
		if len(templateIDs) > 0 {
			var err error
			templates, err = s.api.Template.GetMapByIDs(templateIDs)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		if len(userIDs) > 0 {
			var err error
			users, err = s.api.User.GetMapByIDs(userIDs)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	wg.Wait()

	for _, infoBlock := range infoBlocks {
		if infoBlock.TemplateID != nil {
			infoBlock.Template = templates[*infoBlock.TemplateID]
		}
		if infoBlock.UserID != nil {
			infoBlock.User = users[*infoBlock.UserID]
		}
	}

	return infoBlocks
}

func (s *InfoBlockCollectionService) AggregatesResponses(infoBlocks []*models.InfoBlockResponse) []*models.InfoBlockResponse {
	var templateIDs []uint
	var userIDs []uint
	var infoBlockIDs []uint

	templateIDsMap := make(map[uint]bool)
	userIDsMap := make(map[uint]bool)
	infoBlockIDsMap := make(map[uint]bool)
	infoBlocksInterface := make([]contract.Resource, 0, len(infoBlocks))

	for _, infoBlock := range infoBlocks {
		infoBlocksInterface = append(infoBlocksInterface, infoBlock)

		if infoBlock.TemplateID != nil {
			id := *infoBlock.TemplateID
			if !templateIDsMap[id] {
				templateIDs = append(templateIDs, id)
				templateIDsMap[id] = true
			}
		}

		if infoBlock.UserID != nil {
			id := *infoBlock.UserID
			if !userIDsMap[id] {
				userIDs = append(userIDs, id)
				userIDsMap[id] = true
			}
		}

		id := infoBlock.ID
		if !infoBlockIDsMap[id] {
			infoBlockIDs = append(infoBlockIDs, id)
			infoBlockIDsMap[id] = true
		}
	}

	var wg sync.WaitGroup

	var users map[uint]contract.User
	var templates map[uint]contract.Template
	var galleries map[uuid.UUID][]contract.Gallery
	var allInfoBlocks []*models.InfoBlock

	wg.Add(4)

	go func() {
		defer wg.Done()
		if len(templateIDs) > 0 {
			var err error
			templates, err = s.api.Template.GetMapByIDs(templateIDs)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		if len(userIDs) > 0 {
			var err error
			users, err = s.api.User.GetMapByIDs(userIDs)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		galleries = s.api.Gallery.GetIndexesForResources(infoBlocksInterface)
	}()

	go func() {
		defer wg.Done()
		var err error
		allInfoBlocks, err = s.infoBlockRepo.GetAll()
		if err != nil {
			logger.Error(err)
		}
	}()

	wg.Wait()

	// Создаем карту для быстрого доступа к InfoBlockResponse по ID
	infoBlockMap := make(map[uint]*models.InfoBlockResponse)
	for _, infoBlock := range infoBlocks {
		infoBlockMap[infoBlock.ID] = infoBlock
	}

	// Создаем карту для быстрого доступа к InfoBlock по ID
	allInfoBlockMap := make(map[uint]*models.InfoBlock)
	for _, ib := range allInfoBlocks {
		allInfoBlockMap[ib.ID] = ib
	}

	// Находим детей для каждого InfoBlockResponse
	for _, infoBlock := range infoBlocks {
		// Ищем прямых детей (где InfoBlockID = текущий ID)
		var children []*models.InfoBlockResponse
		for _, ib := range allInfoBlocks {
			if ib.InfoBlockID != nil && *ib.InfoBlockID == infoBlock.ID {
				// Создаем InfoBlockResponse для ребенка
				child := &models.InfoBlockResponse{
					ID:          ib.ID,
					UUID:        ib.UUID,
					TemplateID:  ib.TemplateID,
					UserID:      ib.UserID,
					Media:       ib.Media,
					Title:       ib.Title,
					Description: ib.Description,
					Image:       ib.Image,
				}
				children = append(children, child)
			}
		}
		if len(children) > 0 {
			// Агрегируем детей
			children = s.AggregatesResponses(children)
			infoBlock.Children = children
		}
	}

	// Обогащаем основными данными
	for _, infoBlock := range infoBlocks {
		if gallery, ok := galleries[infoBlock.UUID]; ok {
			infoBlock.Galleries = gallery
		}
		if infoBlock.TemplateID != nil {
			infoBlock.Template = templates[*infoBlock.TemplateID]
		}
		if infoBlock.UserID != nil {
			infoBlock.User = users[*infoBlock.UserID]
		}
	}

	return infoBlocks
}
