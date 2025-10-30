package service

import (
	"github.com/axlle-com/blog/app/models/contract"
	app "github.com/axlle-com/blog/app/service"
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

func (s *MenuItemService) SaveFromRequest(form *request.MenuItemsRequest, user contract.User) (menu *models.MenuItem, err error) {
	menu = app.LoadStruct(&models.MenuItem{}, form).(*models.MenuItem)

	if form.ID == nil {
		err = s.menuItemRepository.Create(menu)
	} else {
		old, err := s.menuItemRepository.GetByID(*form.ID)
		if err != nil {
			return nil, err
		}
		menu.ID = *form.ID
		err = s.menuItemRepository.Update(menu, old)
	}

	if err != nil {
		return
	}

	return
}
