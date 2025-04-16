package service

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/message/models"
	"github.com/axlle-com/blog/pkg/message/repository"
	userProvider "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/google/uuid"
	"sync"
)

type MessageCollectionService struct {
	messageRepo       repository.MessageRepository
	messageService    *MessageService
	userProvider      userProvider.UserProvider
	userGuestProvider userProvider.UserGuestProvider
}

func NewMessageCollectionService(
	messageRepo repository.MessageRepository,
	messageService *MessageService,
	userProvider userProvider.UserProvider,
	userGuestProvider userProvider.UserGuestProvider,
) *MessageCollectionService {
	return &MessageCollectionService{
		messageService:    messageService,
		messageRepo:       messageRepo,
		userProvider:      userProvider,
		userGuestProvider: userGuestProvider,
	}
}

func (s *MessageCollectionService) GetAll() ([]*models.Message, error) {
	return s.messageRepo.GetAll()
}

func (s *MessageCollectionService) GetAllIds() ([]uint, error) {
	return s.messageRepo.GetAllIds()
}

func (s *MessageCollectionService) GetByIDs(ids []uint) ([]*models.Message, error) {
	return s.messageRepo.GetByIDs(ids)
}

func (s *MessageCollectionService) Delete(messages []*models.Message) (err error) {
	var ids []uint
	for _, infoBlock := range messages {
		ids = append(ids, infoBlock.ID)
	}

	if len(ids) > 0 {
		if err = s.messageRepo.DeleteByIDs(ids); err == nil {
			return nil
		}
	}
	return err
}

func (s *MessageCollectionService) WithPaginate(p contracts.Paginator, filter *models.MessageFilter) ([]*models.Message, error) {
	return s.messageRepo.WithPaginate(p, filter)
}

func (s *MessageCollectionService) Aggregates(messages []*models.Message) []*models.Message {
	var userUUIDs []uuid.UUID

	userUUIDsMap := make(map[uuid.UUID]bool)

	for _, message := range messages {
		if message.UserUUID != uuid.Nil {
			if !userUUIDsMap[message.UserUUID] {
				userUUIDs = append(userUUIDs, message.UserUUID)
				userUUIDsMap[message.UserUUID] = true
			}
		}
	}

	var wg sync.WaitGroup

	var users map[uuid.UUID]contracts.User
	var usersGuest map[uuid.UUID]contracts.User

	wg.Add(2)

	go func() {
		defer wg.Done()
		if len(userUUIDs) > 0 {
			var err error
			users, err = s.userProvider.GetMapByUUIDs(userUUIDs)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		if len(userUUIDs) > 0 {
			var err error
			usersGuest, err = s.userGuestProvider.GetMapByUUIDs(userUUIDs)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	wg.Wait()

	for _, message := range messages {
		if message.UserUUID != uuid.Nil {
			if user, ok := users[message.UserUUID]; ok {
				message.User = user
				continue
			}
			message.User = usersGuest[message.UserUUID]

		}
	}

	return messages
}
