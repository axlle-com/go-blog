package service

import (
	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/models/contract"
	selfContracts "github.com/axlle-com/blog/pkg/message/contracts"
	"github.com/axlle-com/blog/pkg/message/form"
	mailer "github.com/axlle-com/blog/pkg/message/job"
)

type MailService struct {
	config                   contract.Config
	queue                    contract.Queue
	messageService           selfContracts.MessageService
	messageCollectionService selfContracts.MessageCollectionService
	api                      *api.Api
}

func NewMailService(
	config contract.Config,
	queue contract.Queue,
	messageService selfContracts.MessageService,
	messageCollectionService selfContracts.MessageCollectionService,
	api *api.Api,

) *MailService {
	return &MailService{
		config:                   config,
		queue:                    queue,
		messageService:           messageService,
		messageCollectionService: messageCollectionService,
		api:                      api,
	}
}

func (s *MailService) SendContact(form *form.Contact) {
	v := s.config.NotifyEmail()
	if v != "" {
		form.To = &v
	}

	s.queue.Enqueue(mailer.NewCreateMessageJob(form), 0)
	s.queue.Enqueue(mailer.NewCreateUserJob(form), 0)
}
