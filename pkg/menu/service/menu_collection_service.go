package service

import (
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/menu/models"
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

func (s *MenuCollectionService) WithPaginate(p contract.Paginator, filter *models.MenuFilter) ([]*models.Menu, error) {
	return s.menuRepository.WithPaginate(p, filter)
}
