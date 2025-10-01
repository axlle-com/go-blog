package service

import (
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/menu/repository"
)

type MenuService struct {
	menuRepository            repository.MenuRepository
	menuItemCollectionService *MenuItemCollectionService
}

func NewMenuService(
	menuRepository repository.MenuRepository,
	menuItemCollectionService *MenuItemCollectionService,
) *MenuService {
	return &MenuService{
		menuRepository:            menuRepository,
		menuItemCollectionService: menuItemCollectionService,
	}
}

func (s *MenuService) GetByID(id uint) (*models.Menu, error) {
	return s.menuRepository.GetByID(id)
}

func (s *MenuService) Aggregate(model *models.Menu) (*models.Menu, error) {
	menuItems, err := s.menuItemCollectionService.GetByParams(map[string]any{"menu_id": model.ID})
	if err != nil {
		return nil, err
	}

	nodes := make(map[uint]*models.MenuItem, len(menuItems))
	roots := make([]*models.MenuItem, 0)

	for _, menuItem := range menuItems {
		nodes[menuItem.ID] = menuItem
	}

	for _, menuItem := range menuItems {
		n := nodes[menuItem.ID]
		if menuItem.MenuItemID == nil {
			roots = append(roots, n)
			continue
		}
		if p, ok := nodes[*menuItem.MenuItemID]; ok {
			p.Children = append(p.Children, n)
		} else {
			roots = append(roots, n)
		}
	}

	model.MenuItems = roots

	return model, err
}
