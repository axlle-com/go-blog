package service

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/menu/repository"
)

type MenuItemCollectionService struct {
	menuItemRepository repository.MenuItemRepository
	menuService        *MenuItemService
}

func NewMenuItemCollectionService(
	menuRepository repository.MenuItemRepository,
	menuService *MenuItemService,
) *MenuItemCollectionService {
	return &MenuItemCollectionService{
		menuItemRepository: menuRepository,
		menuService:        menuService,
	}
}

func (s *MenuItemCollectionService) Filter(paginator contract.Paginator, filter *models.MenuItemFilter) ([]*models.MenuItem, error) {
	collection, err := s.menuItemRepository.GetByFilter(paginator, filter)
	if err != nil {
		return nil, err
	}

	return collection, nil
}

func (s *MenuItemCollectionService) GetByParams(params map[string]any) ([]*models.MenuItem, error) {
	return s.menuItemRepository.GetByParams(params)
}

func (s *MenuItemCollectionService) Aggregate(collection []*models.MenuItem) []*models.MenuItem {
	if len(collection) == 0 {
		return collection
	}

	modelsMap := make(map[uint]*models.MenuItem, len(collection))
	idsMap := make(map[uint][]uint)

	parentSet := make(map[uint]struct{})

	for _, model := range collection {
		if model == nil {
			continue
		}
		modelsMap[model.ID] = model

		if model.MenuItemID != nil && *model.MenuItemID != 0 {
			parentID := *model.MenuItemID
			parentSet[parentID] = struct{}{}
			idsMap[parentID] = append(idsMap[parentID], model.ID)
		}
	}

	if len(parentSet) == 0 {
		return collection
	}

	parentIDs := make([]uint, 0, len(parentSet))
	for id := range parentSet {
		parentIDs = append(parentIDs, id)
	}

	filter := models.NewMenuItemFilter()
	filter.IDs = parentIDs

	parentCollection, err := s.menuItemRepository.GetByFilter(nil, filter)
	if err != nil {
		logger.Errorf("[MenuItemCollectionService][Aggregate] error: %v", err)
		return collection
	}

	parentMap := make(map[uint]*models.MenuItem, len(parentCollection))
	for _, parent := range parentCollection {
		if parent != nil {
			parentMap[parent.ID] = parent
		}
	}

	for parentID, childIDs := range idsMap {
		parent := parentMap[parentID]
		if parent == nil {
			continue
		}
		for _, childID := range childIDs {
			if child := modelsMap[childID]; child != nil {
				child.Parent = parent
			}
		}
	}

	return collection
}

func (s *MenuItemCollectionService) UpdateURLForPublisher(publisher contract.Publisher) (int64, error) {
	return s.menuItemRepository.UpdateURLForPublisher(publisher.GetUUID(), publisher.GetURL())
}

func (s *MenuItemCollectionService) DetachPublisher(publisher contract.Publisher) (int64, error) {
	return s.menuItemRepository.DetachPublisher(publisher.GetUUID())
}
