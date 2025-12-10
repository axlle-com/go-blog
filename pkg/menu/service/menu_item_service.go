package service

import (
	"github.com/axlle-com/blog/app/models/contract"
	app "github.com/axlle-com/blog/app/service/struct"
	"github.com/axlle-com/blog/pkg/menu/http/request"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/menu/repository"
)

type MenuItemService struct {
	menuItemRepository repository.MenuItemRepository
}

func NewMenuItemService(
	menuItemRepository repository.MenuItemRepository,
) *MenuItemService {
	return &MenuItemService{
		menuItemRepository: menuItemRepository,
	}
}

func (s *MenuItemService) SaveFromRequest(form *request.MenuItemsRequest, user contract.User) (menuItem *models.MenuItem, err error) {
	menuItem = app.LoadStruct(&models.MenuItem{}, form).(*models.MenuItem)

	if form.ID == nil {
		err = s.menuItemRepository.Create(menuItem)
	} else {
		old, err := s.menuItemRepository.GetByID(*form.ID)
		if err != nil {
			return nil, err
		}
		menuItem.ID = *form.ID
		err = s.menuItemRepository.Update(menuItem, old)
	}

	return
}

func (s *MenuItemService) GetByID(id uint) (*models.MenuItem, error) {
	return s.menuItemRepository.GetByID(id)
}

func (s *MenuItemService) Delete(id uint) error {
	menuItem, err := s.menuItemRepository.GetByID(id)
	if err != nil {
		return err
	}

	return s.menuItemRepository.Delete(menuItem)
}
