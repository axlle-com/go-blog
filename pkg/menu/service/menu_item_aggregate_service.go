package service

import (
	"net/url"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/menu/repository"
	"github.com/google/uuid"
)

type MenuItemAggregateService struct {
	menuItemRepository repository.MenuItemRepository
	api                *api.Api
}

func NewMenuItemAggregateService(
	menuItemRepository repository.MenuItemRepository,
	api *api.Api,
) *MenuItemAggregateService {
	return &MenuItemAggregateService{
		menuItemRepository: menuItemRepository,
		api:                api,
	}
}

func (s *MenuItemAggregateService) Aggregate(collection []*models.MenuItem) []*models.MenuItem {
	if len(collection) == 0 {
		return collection
	}

	collection = s.enrichWithParents(collection)

	collection = s.enrichWithPublishers(collection)

	return collection
}

func (s *MenuItemAggregateService) enrichWithParents(collection []*models.MenuItem) []*models.MenuItem {
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
		logger.Errorf("[MenuItemAggregateService][enrichWithParents] error: %v", err)
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

func (s *MenuItemAggregateService) enrichWithPublishers(collection []*models.MenuItem) []*models.MenuItem {
	publisherUUIDs := make([]uuid.UUID, 0)
	publisherUUIDsMap := make(map[uuid.UUID]bool)

	for _, item := range collection {
		if item == nil || item.PublisherUUID == nil {
			continue
		}
		u := *item.PublisherUUID
		if !publisherUUIDsMap[u] {
			publisherUUIDs = append(publisherUUIDs, u)
			publisherUUIDsMap[u] = true
		}
	}

	if len(publisherUUIDs) == 0 {
		return collection
	}

	filter := models.NewMenuItemFilter()
	filter.SetUUIDs(publisherUUIDs)

	query := url.Values{}
	query.Set("pageSize", "1000")
	query.Set("page", "1")
	paginator := app.FromQuery(query)

	publishers, _, err := s.api.Publisher.GetPublishers(paginator, filter)
	if err != nil {
		logger.Errorf("[MenuItemAggregateService][enrichWithPublishers] error: %v", err)
		return collection
	}

	publisherMap := make(map[uuid.UUID]contract.Publisher, len(publishers))
	for _, publisher := range publishers {
		if publisher != nil {
			publisherMap[publisher.GetUUID()] = publisher
		}
	}

	for _, item := range collection {
		if item == nil || item.PublisherUUID == nil {
			continue
		}
		if publisher, ok := publisherMap[*item.PublisherUUID]; ok {
			item.Publisher = publisher
		}
	}

	return collection
}
