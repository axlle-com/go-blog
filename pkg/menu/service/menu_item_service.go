package service

import "github.com/axlle-com/blog/pkg/menu/repository"

type MenuItemService struct {
	menuRepository repository.MenuItemRepository
}

func NewMenuItemService(
	menuRepository repository.MenuItemRepository,
) *MenuItemService {
	return &MenuItemService{
		menuRepository: menuRepository,
	}
}
