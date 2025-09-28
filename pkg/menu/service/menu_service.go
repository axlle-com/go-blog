package service

import (
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/menu/repository"
)

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

func (s *MenuService) GetByID(id uint) (*models.Menu, error) {
	return s.menuRepository.GetByID(id)
}
