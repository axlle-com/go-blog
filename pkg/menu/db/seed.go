package db

import (
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	blog "github.com/axlle-com/blog/pkg/blog/provider"
	"github.com/axlle-com/blog/pkg/menu/models"
	"github.com/axlle-com/blog/pkg/menu/repository"
	template "github.com/axlle-com/blog/pkg/template/provider"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"time"
)

type seeder struct {
	menuRepo         repository.MenuRepository
	menuItemRepo     repository.MenuItemRepository
	postProvider     blog.PostProvider
	templateProvider template.TemplateProvider
}

func NewSeeder(
	menu repository.MenuRepository,
	menuItem repository.MenuItemRepository,
	postProvider blog.PostProvider,
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
	err := s.menuItems(n)
	if err != nil {
		return err
	}

	return s.menus(n)
}

func (s *seeder) menus(n int) error {
	ids := s.templateProvider.GetAllIds()
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= n; i++ {
		randomID := ids[rand.Intn(len(ids))]
		_ = models.Menu{
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

		//_, err := s.menuService.Save(&menu, userF)
		//if err != nil {
		//	return err
		//}
	}
	logger.Info("Database seeded Menu successfully!")
	return nil
}

func (s *seeder) menuItems(n int) error {
	rand.Seed(time.Now().UnixNano())

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
	logger.Info("Database seeded Menu successfully!")
	return nil
}
