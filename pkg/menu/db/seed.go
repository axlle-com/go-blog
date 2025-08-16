package db

import (
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

	return s.menuItems(n)
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
			Name:        "Name #" + strconv.Itoa(i),
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

func (s *seeder) menuItems(n int) error {
	for i := 1; i <= n; i++ {
		idsMenu, _ := s.menuRepo.GetAllIds()
		idsMenuItem, _ := s.menuItemRepo.GetAllIds()
		var randomMenuItemID *uint
		if len(idsMenuItem) > 0 {
			randomMenuItemID = &idsMenuItem[rand.Intn(len(idsMenuItem))]
			if rand.Intn(2) == 1 {
				randomMenuItemID = nil
			}
		}

		var randomMenuID uint
		if len(idsMenu) > 0 {
			randomMenuID = idsMenu[rand.Intn(len(idsMenu))]
		}

		menuMenuItem := models.MenuItem{
			MenuID:      randomMenuID,
			MenuItemID:  randomMenuItemID,
			URL:         faker.URL(),
			IsPublished: db.RandBool(),
			Title:       "TitleMenuItem #" + strconv.Itoa(i),
			Sort:        rand.Intn(100),
			CreatedAt:   db.TimePtr(time.Now()),
			UpdatedAt:   db.TimePtr(time.Now()),
			DeletedAt:   nil,
		}

		err := s.menuItemRepo.Create(&menuMenuItem)
		if err != nil {
			return err
		}
	}
	logger.Info("Database seeded MenuItem successfully!")
	return nil
}
