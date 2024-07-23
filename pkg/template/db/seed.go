package db

import (
	. "github.com/axlle-com/blog/pkg/common/db"
	. "github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/template/repository"
	"github.com/bxcodec/faker/v3"
	"log"
	"time"
)

func SeedTemplate(n int) {
	for i := 0; i < n; i++ {
		template := Template{}

		now := time.Now()
		template.Title = faker.Sentence()
		template.Name = faker.Username()
		template.Resource = StrPtr(faker.Username())
		template.CreatedAt = &now
		template.UpdatedAt = &now

		err := repository.NewTemplateRepository().CreateTemplate(&template)
		if err != nil {
			log.Printf("Failed to create template %d: %v", i, err.Error())
		}
	}

	log.Println("Database seeded Template successfully!")
}
