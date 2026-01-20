package db

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/menu/repository"
	publisherModels "github.com/axlle-com/blog/pkg/publisher/models"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
)

type seeder struct {
	config       contract.Config
	disk         contract.DiskService
	seedService  contract.SeedService
	api          *api.Api
	menuRepo     repository.MenuRepository
	menuItemRepo repository.MenuItemRepository
}

type MenuSeedData struct {
	UserID      *uint              `json:"user_id"`
	Title       string             `json:"title"`
	IsMain      *bool              `json:"is_main,omitempty"`
	IsPublished *bool              `json:"is_published,omitempty"`
	Template    *string            `json:"template,omitempty"`
	MenuItems   []MenuItemSeedData `json:"menu_items,omitempty"`
}

type MenuItemSeedData struct {
	PublisherUrl string             `json:"publisher_url"` // Название поста для поиска
	Title        *string            `json:"title,omitempty"`
	URL          *string            `json:"url,omitempty"`
	MenuItemID   *uint              `json:"menu_item_id,omitempty"`
	Sort         *int               `json:"sort,omitempty"`
	Children     []MenuItemSeedData `json:"children,omitempty"`
}

func NewMenuSeeder(
	cfg contract.Config,
	disk contract.DiskService,
	seedService contract.SeedService,
	api *api.Api,
	menu repository.MenuRepository,
	menuItem repository.MenuItemRepository,
) contract.Seeder {
	return &seeder{
		config:       cfg,
		disk:         disk,
		seedService:  seedService,
		api:          api,
		menuRepo:     menu,
		menuItemRepo: menuItem,
	}
}

func (s *seeder) Seed() error {
	return s.seedFromJSON((&models.Menu{}).GetTable())
}

func (s *seeder) seedFromJSON(moduleName string) error {
	files, err := s.seedService.GetFiles(s.config.Layout(), moduleName)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		logger.Infof("[menu][seeder][seedFromJSON] seed files not found for module: %s, skipping", moduleName)
		return nil
	}

	for name, seedPath := range files {
		data, err := s.disk.ReadFile(seedPath)
		if err != nil {
			return err
		}

		ok, err := s.seedService.IsApplied(name)
		if err != nil {
			return err
		}
		if ok {
			continue
		}

		var menusData []MenuSeedData
		if err := json.Unmarshal(data, &menusData); err != nil {
			return err
		}

		for _, menuData := range menusData {
			// Проверяем, существует ли меню с таким title
			found, _ := s.menuRepo.GetByParams(map[string]any{"title": menuData.Title})
			var existingMenu *models.Menu
			if len(found) > 0 {
				existingMenu = found[0]
				logger.Infof("[menu][seeder][seedFromJSON] menu with title='%s' already exists (ID=%d), skipping", menuData.Title, existingMenu.ID)
				continue
			}

			// Создаем меню
			menu := &models.Menu{
				UUID:        uuid.New(),
				Title:       menuData.Title,
				IsPublished: true,
				IsMain:      false,
				Sort:        0,
				CreatedAt:   db.TimePtr(time.Now()),
				UpdatedAt:   db.TimePtr(time.Now()),
			}

			if menuData.IsMain != nil {
				menu.IsMain = *menuData.IsMain
			}
			if menuData.IsPublished != nil {
				menu.IsPublished = *menuData.IsPublished
			}

			// Ищем шаблон, если указан
			if menuData.Template != nil && *menuData.Template != "" {
				resourceName := (&models.Menu{}).GetName()
				tpl, err := s.api.Template.GetByNameAndResource(*menuData.Template, resourceName)
				if err != nil {
					logger.Errorf("[menu][seeder][seedFromJSON] template not found: name=%s, resource=%s, error=%v", *menuData.Template, resourceName, err)
				} else {
					id := tpl.GetID()
					menu.TemplateID = &id
				}
			}

			if err := s.menuRepo.Create(menu); err != nil {
				logger.Errorf("[menu][seeder][seedFromJSON] error creating menu: %v", err)
				continue
			}

			// Создаем пункты меню
			if len(menuData.MenuItems) > 0 {
				sortCounter := 0
				for _, itemData := range menuData.MenuItems {
					if err := s.createMenuItem(menu.ID, nil, itemData, &sortCounter); err != nil {
						logger.Errorf("[menu][seeder][seedFromJSON] error creating menu item: %v", err)
						continue
					}
				}
			}
		}

		if err := s.seedService.MarkApplied(name); err != nil {
			return err
		}

		logger.Infof("[menu][seeder][seedFromJSON] seeded %d menus from JSON (%s)", len(menusData), name)
	}

	return nil
}

func (s *seeder) createMenuItem(menuID uint, parentID *uint, itemData MenuItemSeedData, sortCounter *int) error {
	// Ищем publisher по URL через PublisherProvider
	var publisherUUID *uuid.UUID
	var itemURL string = "#"
	var title string

	if itemData.PublisherUrl != "" {
		// Создаем фильтр по URL
		publisherFilter := publisherModels.NewPublisherFilter()
		publisherFilter.URL = &itemData.PublisherUrl

		// Создаем пагинатор для поиска
		query := make(url.Values)
		query.Set("page", "1")
		query.Set("pageSize", "10")
		paginator := app.FromQuery(query)

		// Ищем publisher через API
		publishers, _, err := s.api.Publisher.GetPublishers(paginator, publisherFilter)
		var foundPublisher contract.Publisher
		if err == nil && len(publishers) > 0 {
			// Ищем точное совпадение по URL
			for _, p := range publishers {
				if p.GetURL() == itemData.PublisherUrl {
					foundPublisher = p
					break
				}
			}
		}

		if foundPublisher != nil {
			newUuid := foundPublisher.GetUUID()
			publisherUUID = &newUuid
			itemURL = foundPublisher.GetURL()
			title = foundPublisher.GetTitle()
		} else {
			logger.Infof("[menu][seeder][createMenuItem] publisher with URL='%s' not found, using as URL", itemData.PublisherUrl)
			itemURL = itemData.PublisherUrl
			title = itemData.PublisherUrl
		}
	}

	if itemData.Title != nil && *itemData.Title != "" {
		title = *itemData.Title
	}

	if title == "" {
		logger.Infof("[menu][seeder][createMenuItem] menu item has no title, skipping")
		return nil
	}

	sort := *sortCounter
	*sortCounter++
	if itemData.Sort != nil {
		sort = *itemData.Sort
	}

	if itemData.URL != nil && *itemData.URL != "" {
		itemURL = *itemData.URL
	}

	// Проверяем, существует ли уже такой пункт меню
	var existingItems []*models.MenuItem
	var err error

	if parentID == nil {
		// Для корневых элементов получаем все пункты меню и фильтруем по menu_item_id IS NULL
		allItems, err := s.menuItemRepo.GetByParams(map[string]any{"menu_id": menuID})
		if err == nil {
			for _, item := range allItems {
				if item.MenuItemID == nil {
					existingItems = append(existingItems, item)
				}
			}
		}
	} else {
		// Для дочерних элементов проверяем по parent_id
		existingItems, err = s.menuItemRepo.GetByParams(map[string]any{
			"menu_id":      menuID,
			"menu_item_id": *parentID,
		})
	}

	if err == nil {
		for _, item := range existingItems {
			// Для корневых элементов проверяем, что menu_item_id действительно NULL
			if parentID == nil && item.MenuItemID != nil {
				continue
			}
			// Для дочерних элементов проверяем, что parent_id совпадает
			if parentID != nil && (item.MenuItemID == nil || *item.MenuItemID != *parentID) {
				continue
			}

			// Проверяем, совпадает ли URL или publisher_uuid
			if publisherUUID != nil && item.PublisherUUID != nil && *item.PublisherUUID == *publisherUUID {
				logger.Infof("[menu][seeder][createMenuItem] menu item already exists (menu_id=%d, parent_id=%v, publisher_uuid='%s'), skipping", menuID, parentID, *publisherUUID)
				// Создаем дочерние пункты для существующего элемента
				for _, childData := range itemData.Children {
					if err := s.createMenuItem(menuID, &item.ID, childData, sortCounter); err != nil {
						logger.Errorf("[menu][seeder][createMenuItem] error creating child menu item: %v", err)
						continue
					}
				}
				return nil
			}
			if publisherUUID == nil && item.PublisherUUID == nil && item.URL == itemURL {
				logger.Infof("[menu][seeder][createMenuItem] menu item already exists (menu_id=%d, parent_id=%v, url='%s'), skipping", menuID, parentID, itemURL)
				// Создаем дочерние пункты для существующего элемента
				for _, childData := range itemData.Children {
					if err := s.createMenuItem(menuID, &item.ID, childData, sortCounter); err != nil {
						logger.Errorf("[menu][seeder][createMenuItem] error creating child menu item: %v", err)
						continue
					}
				}
				return nil
			}
		}
	}

	menuItem := &models.MenuItem{
		MenuID:        menuID,
		MenuItemID:    parentID,
		PublisherUUID: publisherUUID,
		Title:         title,
		URL:           itemURL,
		Sort:          sort,
		CreatedAt:     db.TimePtr(time.Now()),
		UpdatedAt:     db.TimePtr(time.Now()),
	}

	if err := s.menuItemRepo.Create(menuItem); err != nil {
		return err
	}

	// Создаем дочерние пункты
	for _, childData := range itemData.Children {
		if err := s.createMenuItem(menuID, &menuItem.ID, childData, sortCounter); err != nil {
			logger.Errorf("[menu][seeder][createMenuItem] error creating child menu item: %v", err)
			continue
		}
	}

	return nil
}

func (s *seeder) SeedTest(n int) error {
	err := s.menus(n)
	if err != nil {
		return err
	}

	return s.seedItemsPerMenu()
}

func (s *seeder) menus(n int) error {
	ids := s.api.Template.GetAllIds()
	for i := 1; i <= n; i++ {
		randomID := ids[rand.Intn(len(ids))]
		menu := &models.Menu{
			UUID:        uuid.New(),
			TemplateID:  &randomID,
			IsPublished: db.RandBool(),
			IsMain:      db.RandBool(),
			Title:       "Name #" + strconv.Itoa(i),
			Ico:         db.StrPtr("Ico #" + strconv.Itoa(i)),
			Sort:        rand.Intn(100),
			CreatedAt:   db.TimePtr(time.Now()),
			UpdatedAt:   db.TimePtr(time.Now()),
			DeletedAt:   nil,
		}

		err := s.menuRepo.Create(menu)
		if err != nil {
			return err
		}
	}
	logger.Info("Database seeded Menu successfully!")
	return nil
}

func (s *seeder) seedItemsPerMenu() error {
	menuIDs, err := s.menuRepo.GetAllIds()
	if err != nil {
		return err
	}
	if len(menuIDs) == 0 {
		return fmt.Errorf("no menus found")
	}

	for _, menuID := range menuIDs {
		if err := s.seedTreeForMenu(menuID); err != nil {
			return err
		}
	}
	logger.Info("Database seeded MenuItem successfully!")
	return nil
}

func (s *seeder) seedTreeForMenu(menuID uint) error {
	now := time.Now()

	// локальный helper: создаёт пункт в рамках menuID; parentID — внутри того же меню
	createItem := func(parentID *uint, sort int, title string) (*models.MenuItem, error) {
		it := &models.MenuItem{
			MenuID:     menuID,
			MenuItemID: parentID,
			URL:        faker.URL(),
			Title:      title,
			Sort:       sort,
			CreatedAt:  db.TimePtr(now),
			UpdatedAt:  db.TimePtr(now),
		}
		if err := s.menuItemRepo.Create(it); err != nil {
			return nil, err
		}
		// repo.Create сам выставит Path:
		//  - корень: "/<id>/"
		//  - дочерний: "<parent.Path><id>/"
		return it, nil
	}

	// 1) корни
	rootCount := randIntRange(3, 6)
	roots := make([]*models.MenuItem, 0, rootCount)
	for i := 0; i < rootCount; i++ {
		item, err := createItem(nil, i*10, fmt.Sprintf("Root %d", i+1))
		if err != nil {
			return err
		}
		roots = append(roots, item)
	}

	// 2) дети и 3) внуки (всё в том же menu_id)
	for _, r := range roots {
		childCount := randIntRange(0, 4)
		children := make([]*models.MenuItem, 0, childCount)
		for j := 0; j < childCount; j++ {
			item, err := createItem(&r.ID, j*10, fmt.Sprintf("Child %d of %d", j+1, r.ID))
			if err != nil {
				return err
			}
			children = append(children, item)
		}
		for _, c := range children {
			grandCount := randIntRange(0, 3)
			for k := 0; k < grandCount; k++ {
				if _, err := createItem(&c.ID, k*10, fmt.Sprintf("Grand %d of %d", k+1, c.ID)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func randIntRange(min, max int) int {
	if max < min {
		min, max = max, min
	}
	return rand.Intn(max-min+1) + min
}
