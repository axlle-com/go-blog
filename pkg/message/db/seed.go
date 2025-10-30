package db

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/message/models"
	"github.com/axlle-com/blog/pkg/message/service"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/bxcodec/faker/v3"
)

type seeder struct {
	messageService *service.MessageService
	userProvider   user.UserProvider
}

func NewMessageSeeder(
	messageService *service.MessageService,
	userProvider user.UserProvider,
) contract.Seeder {
	return &seeder{
		messageService: messageService,
		userProvider:   userProvider,
	}
}

func (s *seeder) Seed() error {
	return nil
}

func (s *seeder) SeedTest(n int) error {
	idsUser := s.userProvider.GetAll()
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
