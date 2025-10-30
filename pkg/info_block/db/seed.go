package db

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/bxcodec/faker/v3"
)

type seeder struct {
	infoBlockService *service.InfoBlockService
	templateProvider template.TemplateProvider
	userProvider     user.UserProvider
}

func NewSeeder(
	infoBlockService *service.InfoBlockService,
	templateProvider template.TemplateProvider,
	user user.UserProvider,
) contract.Seeder {
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

	for i := 1; i <= n; i++ {
		randomID := ids[rand.Intn(len(ids))]
		randomUserID := idsUser[rand.Intn(len(idsUser))]

		infoBlock := models.InfoBlock{
			TemplateID:  &randomID,
			Media:       db.StrPtr(faker.Word()),
			Title:       "TitleInfoBlock #" + strconv.Itoa(i),
			Description: db.StrPtr(faker.Paragraph()),
			Image:       db.StrPtr("/public/img/404.svg"),
			CreatedAt:   db.TimePtr(time.Now()),
			UpdatedAt:   db.TimePtr(time.Now()),
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
