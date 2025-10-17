package service

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/models/contracts"
	app "github.com/axlle-com/blog/app/service"
	"github.com/axlle-com/blog/pkg/menu/http/request"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/menu/repository"
)

type MenuService struct {
	menuRepository            repository.MenuRepository
	menuItemService           *MenuItemService
	menuItemCollectionService *MenuItemCollectionService
}

func NewMenuService(
	menuRepository repository.MenuRepository,
	menuItemService *MenuItemService,
	menuItemCollectionService *MenuItemCollectionService,
) *MenuService {
	return &MenuService{
		menuRepository:            menuRepository,
		menuItemService:           menuItemService,
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

	menuItems = s.menuItemCollectionService.Aggregate(menuItems)

	nodes := make(map[uint]*models.MenuItem, len(menuItems))
	roots := make([]*models.MenuItem, 0)

	for _, menuItem := range menuItems {
		nodes[menuItem.ID] = menuItem
	}

	for _, menuItem := range menuItems {
		item := nodes[menuItem.ID]
		if menuItem.MenuItemID == nil {
			roots = append(roots, item)
			continue
		}
		if itemInNode, ok := nodes[*menuItem.MenuItemID]; ok {
			itemInNode.Children = append(itemInNode.Children, item)
		} else {
			roots = append(roots, item)
		}
	}

	model.MenuItems = roots

	return model, err
}

func (s *MenuService) SaveFromRequest(form *request.MenuRequest, found *models.Menu, user contracts.User) (*models.Menu, error) {
	newMenu := app.LoadStruct(&models.Menu{}, form).(*models.Menu)

	// 1) создаём/обновляем меню
	if found == nil {
		if err := s.menuRepository.Create(newMenu); err != nil {
			return nil, err
		}
	} else {
		newMenu.ID = found.ID
		if err := s.menuRepository.Update(newMenu); err != nil {
			return nil, err
		}
	}

	itemsReq := form.MenuItems
	if len(itemsReq) == 0 {
		return newMenu, nil
	}

	// 2) параллельно сохраняем пункты меню
	results := make([]*models.MenuItem, len(itemsReq))

	// ограничим количество одновременно работающих горутин
	maxPar := 2 * runtime.GOMAXPROCS(0)
	if maxPar < 4 {
		maxPar = 4
	}
	if maxPar > len(itemsReq) {
		maxPar = len(itemsReq)
	}
	sem := make(chan struct{}, maxPar)

	var wg sync.WaitGroup

	newErr := errutil.New()

	for i := range itemsReq {
		if itemsReq[i] == nil {
			continue
		}

		i := i // захват индекса!
		sem <- struct{}{}

		app.SafeGo(&wg, func() {
			defer func() { <-sem }()

			item, err := s.menuItemService.SaveFromRequest(itemsReq[i], user)
			if err != nil {
				newErr.Add(fmt.Errorf("[MenuItemService][SaveFromRequest] err: %s", err.Error()))
				return
			}

			if item == nil {
				newErr.Add(fmt.Errorf("[MenuItemService][SaveFromRequest] item is nil"))
				return
			}

			// писать в заранее выделенный слайс безопасно по разным индексам
			results[i] = item
		})
	}

	wg.Wait()

	compact := results[:0]
	for _, it := range results {
		if it != nil {
			compact = append(compact, it)
		}
	}
	newMenu.MenuItems = compact

	return newMenu, newErr.Error()
}
