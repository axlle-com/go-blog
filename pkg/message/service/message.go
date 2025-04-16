package service

import (
	"github.com/axlle-com/blog/app/models/contracts"
	app "github.com/axlle-com/blog/app/service"
	"github.com/axlle-com/blog/pkg/message/models"
	"github.com/axlle-com/blog/pkg/message/repository"
	userProvider "github.com/axlle-com/blog/pkg/user/provider"
)

type MessageService struct {
	messageRepo       repository.MessageRepository
	userProvider      userProvider.UserProvider
	userGuestProvider userProvider.UserGuestProvider
}

func NewMessageService(
	messageRepository repository.MessageRepository,
	userProvider userProvider.UserProvider,
	userGuestProvider userProvider.UserGuestProvider,
) *MessageService {
	return &MessageService{
		messageRepo:       messageRepository,
		userProvider:      userProvider,
		userGuestProvider: userGuestProvider,
	}
}

func (s *MessageService) GetByID(id uint) (*models.Message, error) {
	return s.messageRepo.GetByID(id)
}

func (s *MessageService) Aggregate(message *models.Message) *models.Message {
	return message
}

func (s *MessageService) Create(message *models.Message, user contracts.User) (*models.Message, error) {
	if user != nil {
		message.UserUUID = user.GetUUID()
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

func (s *MessageService) SaveFromRequest(form *models.MessageRequest, found *models.Message, user contracts.User) (message *models.Message, err error) {
	templateForm := app.LoadStruct(&models.Message{}, form).(*models.Message)

	if found == nil {
		message, err = s.Create(templateForm, user)
	} else {
		templateForm.ID = found.ID
		message, err = s.Update(templateForm)
	}

	if err != nil {
		return
	}

	return
}
