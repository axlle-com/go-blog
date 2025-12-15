package db

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/message/models"
	"github.com/axlle-com/blog/pkg/message/service"
	"github.com/bxcodec/faker/v3"
)

type seeder struct {
	messageService *service.MessageService
	api            *api.Api
}

func NewMessageSeeder(
	messageService *service.MessageService,
	api *api.Api,
) contract.Seeder {
	return &seeder{
		messageService: messageService,
		api:            api,
	}
}

func (s *seeder) Seed() error {
	return nil
}

func (s *seeder) SeedTest(n int) error {
	idsUser := s.api.User.GetAll()
	for i := 1; i <= n; i++ {
		randomUser := idsUser[rand.Intn(len(idsUser))]

		message := &models.Message{}

		now := time.Now()
		message.Subject = db.StrPtr("GetSubject #" + strconv.Itoa(i))
		message.To = db.StrPtr(strconv.Itoa(i) + "_to@mail.com")
		message.Viewed = false
		message.From = db.StrPtr(strconv.Itoa(i) + "_from@mail.com")
		message.Body = faker.Paragraph()
		message.CreatedAt = &now
		message.UpdatedAt = &now

		_, err := s.messageService.Create(message, randomUser.GetUUID().String())
		if err != nil {
			return err
		}
	}

	logger.Info("Database seeded Message successfully!")
	return nil
}
