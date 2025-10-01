package db

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/menu/repository"
	template "github.com/axlle-com/blog/pkg/template/provider"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
)

type seeder struct {
	menuRepo         repository.MenuRepository
	menuItemRepo     repository.MenuItemRepository
	postProvider     contracts.PostProvider
	templateProvider template.TemplateProvider
}

func NewMenuSeeder(
	menu repository.MenuRepository,
	menuItem repository.MenuItemRepository,
	postProvider contracts.PostProvider,
	template template.TemplateProvider,
) contracts.Seeder {
	return &seeder{
		menuRepo:         menu,
		menuItemRepo:     menuItem,
		postProvider:     postProvider,
		templateProvider: template,
	}
}

func (s *seeder) Seed() error {
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
	ids := s.templateProvider.GetAllIds()
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
			MenuID:      menuID,
			MenuItemID:  parentID,
			URL:         faker.URL(),
			IsPublished: db.RandBool(),
			Title:       title,
			Sort:        sort,
			CreatedAt:   db.TimePtr(now),
			UpdatedAt:   db.TimePtr(now),
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
