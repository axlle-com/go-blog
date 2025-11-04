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

	templateIDsMap := make(map[uint]bool)
	userIDsMap := make(map[uint]bool)
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
	}

	var wg sync.WaitGroup

	var users map[uint]contract.User
	var templates map[uint]contract.Template
	var galleries map[uuid.UUID][]contract.Gallery

	wg.Add(3)

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

	wg.Wait()

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
