package service

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
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

func (s *MenuItemService) SaveFromRequest(form *request.MenuItemsRequest, user contracts.User) (menu *models.MenuItem, err error) {
	newMenu := app.LoadStruct(&models.MenuItem{}, form).(*models.MenuItem)

	logger.Dump(newMenu)
	if form.ID == nil {
		err = s.menuItemRepository.Create(newMenu)
	} else {
		old, err := s.menuItemRepository.GetByID(*form.ID)
		if err != nil {
			return nil, err
		}
		newMenu.ID = *form.ID
		err = s.menuItemRepository.Update(newMenu, old)
	}

	if err != nil {
		return
	}

	return
}
