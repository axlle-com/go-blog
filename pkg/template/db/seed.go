package db

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/template/models"
	"github.com/axlle-com/blog/pkg/template/repository"
	"github.com/bxcodec/faker/v3"
)

type seeder struct {
	template repository.TemplateRepository
}

func NewSeeder(
	template repository.TemplateRepository,
) contracts.Seeder {
	return &seeder{
		template: template,
	}
}

func (s *seeder) Seed() error {
	return nil
}

func (s *seeder) SeedTest(n int) error {
	resources := map[string]string{
		"Шаблон для постов":    "posts",
		"Шаблон для категорий": "post_categories",
		"Шаблон для тегов":     "post_tags",
		"Шаблон для блоков":    "info_blocks",
		"Шаблон для меню":      "menus",
	}

	// один раз посеять rng
	rand.Seed(time.Now().UnixNano())

	// подготовим слайс ключей
	keys := make([]string, 0, len(resources))
	for title := range resources {
		keys = append(keys, title)
	}

	for i := 1; i <= n; i++ {
		var template models.Template
		now := time.Now()

		// берём случайный ключ и значение из map
		title := keys[rand.Intn(len(keys))]
		resource := resources[title]

		template.Title = fmt.Sprintf("%s #%d", title, i)
		template.Name = faker.Username()
		template.ResourceName = db.StrPtr(resource)
		template.CreatedAt = &now
		template.UpdatedAt = &now

		if err := s.template.Create(&template); err != nil {
			return err
		}
	}

	logger.Info("Database seeded Template successfully!")
	return nil
}
