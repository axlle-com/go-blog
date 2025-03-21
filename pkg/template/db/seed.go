package db

import (
	"github.com/axlle-com/blog/pkg/app/db"
	"github.com/axlle-com/blog/pkg/app/logger"
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	. "github.com/axlle-com/blog/pkg/template/models"
	"github.com/axlle-com/blog/pkg/template/repository"
	"github.com/bxcodec/faker/v3"
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

func (s *seeder) Seed() {}

func (s *seeder) SeedTest(n int) {
	for i := 0; i < n; i++ {
		template := Template{}

		now := time.Now()
		template.Title = faker.Sentence()
		template.Name = faker.Username()
		template.Tabular = db.StrPtr(faker.Username())
		template.CreatedAt = &now
		template.UpdatedAt = &now

		err := s.template.Create(&template)
		if err != nil {
			logger.Errorf("Failed to create template %d: %v", i, err.Error())
		}
	}

	logger.Info("Database seeded Template successfully!")
}

func SeedTemplate(n int) {

}
