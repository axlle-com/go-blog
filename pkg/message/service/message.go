package service

import (
	"github.com/axlle-com/blog/app/logger"
	appContracts "github.com/axlle-com/blog/app/models/contract"
	app "github.com/axlle-com/blog/app/service"
	"github.com/axlle-com/blog/pkg/message/contracts"
	"github.com/axlle-com/blog/pkg/message/models"
	userProvider "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/google/uuid"
)

type MessageService struct {
	messageRepo  contracts.MessageRepository
	userProvider userProvider.UserProvider
}

func NewMessageService(
	messageRepository contracts.MessageRepository,
	userProvider userProvider.UserProvider,
) *MessageService {
	return &MessageService{
		messageRepo:  messageRepository,
		userProvider: userProvider,
	}
}

func (s *MessageService) GetByID(id uint) (*models.Message, error) {
	return s.messageRepo.GetByID(id)
}

func (s *MessageService) Aggregate(message *models.Message) *models.Message {
	var user appContracts.User

	if message.UserUUID != uuid.Nil {
		var err error
		user, err = s.userProvider.GetByUUID(message.UserUUID)
		if err != nil {
			logger.Errorf("[MessageService][Aggregates] Error: %v", err)
		}
	}

	message.User = user

	return message
}

func (s *MessageService) Create(message *models.Message, userUuid string) (*models.Message, error) {
	if userUuid != "" {
		newUUID, err := uuid.Parse(userUuid)
		if err != nil {
			logger.Errorf("Invalid UUID: %v", err)
		}
		message.UserUUID = newUUID
	}
	if err := s.messageRepo.Create(message); err != nil {
		return nil, err
	}
	return message, nil
}

func (s *MessageService) Update(message *models.Message) (*models.Message, error) {
	if err := s.messageRepo.Update(message); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *MessageService) Delete(message *models.Message) (err error) {
	return s.messageRepo.Delete(message)
}

func (s *MessageService) SaveFromRequest(form *models.MessageRequest, found *models.Message, user appContracts.User) (message *models.Message, err error) {
	templateForm := app.LoadStruct(&models.Message{}, form).(*models.Message)

	if found == nil {
		message, err = s.Create(templateForm, user.GetUUID().String())
	} else {
		templateForm.ID = found.ID
		message, err = s.Update(templateForm)
	}

	if err != nil {
		return
	}

	return
}
