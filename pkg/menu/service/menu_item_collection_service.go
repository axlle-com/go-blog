package service

import (
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

func (s *MenuItemCollectionService) UpdateURLForPublisher(publisher contract.Publisher) (int64, error) {
	return s.menuItemRepository.UpdateURLForPublisher(publisher.GetUUID(), publisher.GetURL())
}

func (s *MenuItemCollectionService) DetachPublisher(publisher contract.Publisher) (int64, error) {
	return s.menuItemRepository.DetachPublisher(publisher.GetUUID())
}
