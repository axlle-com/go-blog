package service

import (
	"github.com/axlle-com/blog/pkg/menu/repository"
)

type MenuCollectionService struct {
	menuRepository repository.MenuRepository
	menuService    *MenuService
}

func NewMenuCollectionService(
	menuRepository repository.MenuRepository,
	menuService *MenuService,
) *MenuCollectionService {
	return &MenuCollectionService{
		menuRepository: menuRepository,
		menuService:    menuService,
	}
}
