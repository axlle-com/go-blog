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

func (s *seeder) Seed() error {
	return nil
}

func (s *seeder) SeedTest(n int) error {
	return s.infoBlocks(n)
}

func (s *seeder) infoBlocks(n int) error {
	idsUser := s.userProvider.GetAllIds()
	ids := s.templateProvider.GetAllIds()
	rand.Seed(time.Now().UnixNano())

	for i := 1; i <= n; i++ {
		randomID := ids[rand.Intn(len(ids))]
		randomUserID := idsUser[rand.Intn(len(idsUser))]

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
			return err
		}
	}
	logger.Info("Database seeded infoBlock successfully!")
	return nil
}
