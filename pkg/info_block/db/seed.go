package db

import (
	. "github.com/axlle-com/blog/pkg/app/db"
	"github.com/axlle-com/blog/pkg/app/logger"
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	. "github.com/axlle-com/blog/pkg/info_block/models"
	. "github.com/axlle-com/blog/pkg/info_block/repository"
	. "github.com/axlle-com/blog/pkg/info_block/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	"github.com/bxcodec/faker/v3"
	"math/rand"
	"time"
)

type seeder struct {
	infoBlockRepo    InfoBlockRepository
	infoBlockService *InfoBlockService
	templateProvider template.TemplateProvider
}

func NewSeeder(
	infoBlock InfoBlockRepository,
	infoBlockService *InfoBlockService,
	templateProvider template.TemplateProvider,
) contracts.Seeder {
	return &seeder{
		infoBlockRepo:    infoBlock,
		infoBlockService: infoBlockService,
		templateProvider: templateProvider,
	}
}

func (s *seeder) Seed() {}

func (s *seeder) SeedTest(n int) {
	s.infoBlocks(n)
}

func (s *seeder) infoBlocks(n int) {
	ids := s.templateProvider.GetAllIds()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		randomID := ids[rand.Intn(len(ids))]
		infoBlock := InfoBlock{
			TemplateID:  &randomID,
			Media:       StrPtr(faker.Word()),
			Title:       faker.Sentence(),
			Description: StrPtr(faker.Paragraph()),
			Image:       StrPtr("/public/img/404.svg"),
			CreatedAt:   TimePtr(time.Now()),
			UpdatedAt:   TimePtr(time.Now()),
			DeletedAt:   nil,
		}

		err := s.infoBlockRepo.Create(&infoBlock)
		if err != nil {
			logger.Errorf("Failed to create infoBlock %d: %v", i, err.Error())
		}
	}
	logger.Info("Database seeded infoBlock successfully!")
}
