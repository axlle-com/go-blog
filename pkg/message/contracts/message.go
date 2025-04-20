package contracts

import (
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/message/models"
)

type MessageService interface {
	GetByID(id uint) (*models.Message, error)
	Aggregate(message *models.Message) *models.Message
	Create(message *models.Message, userUuid string) (*models.Message, error)
	Update(message *models.Message) (*models.Message, error)
	Delete(message *models.Message) (err error)
	SaveFromRequest(form *models.MessageRequest, found *models.Message, user contracts.User) (message *models.Message, err error)
}
