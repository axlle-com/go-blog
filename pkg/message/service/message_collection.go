package service

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	contracts2 "github.com/axlle-com/blog/pkg/message/contracts"
	"github.com/axlle-com/blog/pkg/message/models"
	userProvider "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/google/uuid"
)

type MessageCollectionService struct {
	messageRepo    contracts2.MessageRepository
	messageService contracts2.MessageService
	userProvider   userProvider.UserProvider
}

func NewMessageCollectionService(
	messageRepo contracts2.MessageRepository,
	messageService contracts2.MessageService,
	userProvider userProvider.UserProvider,
) *MessageCollectionService {
	return &MessageCollectionService{
		messageService: messageService,
		messageRepo:    messageRepo,
		userProvider:   userProvider,
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

func (s *MessageCollectionService) CountByField(field string, value any) (int64, error) {
	return s.messageRepo.CountByField(field, value)
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

func (s *MessageCollectionService) Paginator(p contracts.Paginator, filter *models.MessageFilter) (contracts.Paginator, error) {
	return s.messageRepo.Paginator(p, filter)
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

	var users map[uuid.UUID]contracts.User

	if len(userUUIDs) > 0 {
		var err error
		users, err = s.userProvider.GetMapByUUIDs(userUUIDs)
		if err != nil {
			logger.Errorf("[MessageCollectionService][Aggregates] Error: %v", err)
		}
	}

	for _, message := range messages {
		if message.UserUUID != uuid.Nil {
			if user, ok := users[message.UserUUID]; ok {
				message.User = user
			}

		}
	}

	return messages
}
