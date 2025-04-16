package db

import (
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/message/models"
	"github.com/axlle-com/blog/pkg/message/service"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/bxcodec/faker/v3"
	"math/rand"
	"strconv"
	"time"
)

type seeder struct {
	messageService *service.MessageService
	userProvider   user.UserProvider
}

func NewMessageSeeder(
	messageService *service.MessageService,
	userProvider user.UserProvider,
) contracts.Seeder {
	return &seeder{
		messageService: messageService,
		userProvider:   userProvider,
	}
}

func (s *seeder) Seed() {}

func (s *seeder) SeedTest(n int) {
	idsUser := s.userProvider.GetAll()
	for i := 1; i <= n; i++ {
		randomUser := idsUser[rand.Intn(len(idsUser))]

		message := &models.Message{}
		now := time.Now()
		message.Subject = db.StrPtr("Subject #" + strconv.Itoa(i))
		message.Body = faker.Username()
		message.CreatedAt = &now
		message.UpdatedAt = &now

		_, err := s.messageService.Create(message, randomUser)
		if err != nil {
			logger.Errorf("Failed to create message %d: %v", i, err.Error())
		}
	}

	logger.Info("Database seeded Message successfully!")
}
