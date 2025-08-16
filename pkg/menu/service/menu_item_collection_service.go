package service

import (
	"github.com/axlle-com/blog/pkg/menu/repository"
)

type MenuItemCollectionService struct {
	menuRepository repository.MenuItemRepository
	menuService    *MenuItemService
}

func NewMenuItemCollectionService(
	menuRepository repository.MenuItemRepository,
	menuService *MenuItemService,
) *MenuItemCollectionService {
	return &MenuItemCollectionService{
		menuRepository: menuRepository,
		menuService:    menuService,
	}
}
