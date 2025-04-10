package db

import (
	. "github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	. "github.com/axlle-com/blog/pkg/info_block/models"
	. "github.com/axlle-com/blog/pkg/info_block/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/bxcodec/faker/v3"
	"math/rand"
	"strconv"
	"time"
)

type seeder struct {
	infoBlockService *InfoBlockService
	templateProvider template.TemplateProvider
	userProvider     user.UserProvider
}

func NewSeeder(
	infoBlockService *InfoBlockService,
	templateProvider template.TemplateProvider,
	user user.UserProvider,
) contracts.Seeder {
	return &seeder{
		infoBlockService: infoBlockService,
		templateProvider: templateProvider,
		userProvider:     user,
	}
}

func (s *seeder) Seed() {}

func (s *seeder) SeedTest(n int) {
	s.infoBlocks(n)
}

func (s *seeder) infoBlocks(n int) {
	idsUser := s.userProvider.GetAllIds()
	randomUserID := idsUser[rand.Intn(len(idsUser))]
	ids := s.templateProvider.GetAllIds()
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= n; i++ {
		randomID := ids[rand.Intn(len(ids))]
		infoBlock := InfoBlock{
			TemplateID:  &randomID,
			Media:       StrPtr(faker.Word()),
			Title:       "TitleInfoBlock #" + strconv.Itoa(i),
			Description: StrPtr(faker.Paragraph()),
			Image:       StrPtr("/public/img/404.svg"),
			CreatedAt:   TimePtr(time.Now()),
			UpdatedAt:   TimePtr(time.Now()),
			DeletedAt:   nil,
			UserID:      &randomUserID,
		}

		_, err := s.infoBlockService.Create(&infoBlock, nil)
		if err != nil {
			logger.Errorf("Failed to create infoBlock %d: %v", i, err.Error())
		}
	}
	logger.Info("Database seeded infoBlock successfully!")
}
