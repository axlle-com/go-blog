package service

import "github.com/axlle-com/blog/pkg/menu/repository"

type MenuService struct {
	menuRepository repository.MenuRepository
}

func NewMenuService(
	menuRepository repository.MenuRepository,
) *MenuService {
	return &MenuService{
		menuRepository: menuRepository,
	}
}
