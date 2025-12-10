package service

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/models/contract"
	app "github.com/axlle-com/blog/app/service"
	"github.com/axlle-com/blog/app/service/struct"
	"github.com/axlle-com/blog/pkg/menu/http/request"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/menu/repository"
)

type MenuService struct {
	menuRepo                  repository.MenuRepository
	menuItemService           *MenuItemService
	menuItemCollectionService *MenuItemCollectionService
	menuItemAggregateService  *MenuItemAggregateService
}

func NewMenuService(
	menuRepository repository.MenuRepository,
	menuItemService *MenuItemService,
	menuItemCollectionService *MenuItemCollectionService,
	menuItemAggregateService *MenuItemAggregateService,
) *MenuService {
	return &MenuService{
		menuRepo:                  menuRepository,
		menuItemService:           menuItemService,
		menuItemCollectionService: menuItemCollectionService,
		menuItemAggregateService:  menuItemAggregateService,
	}
}

func (s *MenuService) GetByID(id uint) (*models.Menu, error) {
	return s.menuRepo.GetByID(id)
}

func (s *MenuService) Aggregate(model *models.Menu) (*models.Menu, error) {
	menuItems, err := s.menuItemCollectionService.GetByParams(map[string]any{"menu_id": model.ID})
	if err != nil {
		return nil, err
	}

	menuItems = s.menuItemAggregateService.Aggregate(menuItems)

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

func (s *MenuService) SaveFromRequest(form *request.MenuRequest, found *models.Menu, user contract.User) (*models.Menu, error) {
	newMenu := _struct.LoadStruct(&models.Menu{}, form).(*models.Menu)

	// 1) создаём/обновляем меню
	if found == nil {
		if err := s.menuRepo.Create(newMenu); err != nil {
			return nil, err
		}
	} else {
		newMenu.ID = found.ID
		if err := s.menuRepo.Update(newMenu); err != nil {
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

		itemsReq[i].MenuID = newMenu.ID

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

func (s *MenuService) GetMenuWithItems(menuID uint) (*models.Menu, error) {
	menu, err := s.menuRepo.GetByID(menuID)
	if err != nil {
		return nil, fmt.Errorf("menu not found: %w", err)
	}

	filter := &models.MenuItemFilter{
		MenuID: &menuID,
	}

	// paginator = nil → вернёт всё
	items, err := s.menuItemCollectionService.Filter(nil, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to load menu items: %w", err)
	}

	roots := buildMenuTree(items)
	menu.MenuItems = roots

	return menu, nil
}

func buildMenuTree(items []*models.MenuItem) []*models.MenuItem {
	byID := make(map[uint]*models.MenuItem, len(items))
	var roots []*models.MenuItem

	// подготовим карту по ID
	for _, item := range items {
		item.Children = nil
		item.Parent = nil
		byID[item.ID] = item
	}

	// склеиваем дерево
	for _, item := range items {
		if item.MenuItemID == nil || *item.MenuItemID == 0 {
			roots = append(roots, item)
			continue
		}

		parent, ok := byID[*item.MenuItemID]
		if !ok {
			// на всякий случай не теряем узел
			roots = append(roots, item)
			continue
		}

		parent.Children = append(parent.Children, item)
		item.Parent = parent
	}

	return roots
}
