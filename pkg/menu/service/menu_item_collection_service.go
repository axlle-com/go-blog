package service

import (
	"github.com/axlle-com/blog/app/models/contracts"
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

func (s *MenuItemCollectionService) Filter(paginator contracts.Paginator, filter *models.MenuItemFilter) ([]*models.MenuItem, error) {
	return s.menuItemRepository.GetByFilter(paginator, filter)
}

func (s *MenuItemCollectionService) GetByParams(params map[string]any) ([]*models.MenuItem, error) {
	return s.menuItemRepository.GetByParams(params)
}
