package db

import (
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	. "github.com/axlle-com/blog/pkg/template/models"
	"github.com/axlle-com/blog/pkg/template/repository"
	"github.com/bxcodec/faker/v3"
	"strconv"
	"time"
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
	for i := 1; i <= n; i++ {
		template := Template{}

		now := time.Now()
		template.Title = "TitleTemplate #" + strconv.Itoa(i)
		template.Name = faker.Username()
		template.ResourceName = db.StrPtr(faker.Username())
		template.CreatedAt = &now
		template.UpdatedAt = &now

		err := s.template.Create(&template)
		if err != nil {
			return err
		}
	}

	logger.Info("Database seeded Template successfully!")
	return nil
}
