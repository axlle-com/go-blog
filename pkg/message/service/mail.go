package service

import (
	"github.com/axlle-com/blog/app/models/contracts"
	selfContracts "github.com/axlle-com/blog/pkg/message/contracts"
	"github.com/axlle-com/blog/pkg/message/form"
	mailer "github.com/axlle-com/blog/pkg/message/job"
	"github.com/axlle-com/blog/pkg/user/provider"
)

type MailService struct {
	config                   contracts.Config
	queue                    contracts.Queue
	messageService           selfContracts.MessageService
	messageCollectionService selfContracts.MessageCollectionService
	userProvider             provider.UserProvider
}

func NewMailService(
	config contracts.Config,
	queue contracts.Queue,
	messageService selfContracts.MessageService,
	messageCollectionService selfContracts.MessageCollectionService,
	userProvider provider.UserProvider,

) *MailService {
	return &MailService{
		config:                   config,
		queue:                    queue,
		messageService:           messageService,
		messageCollectionService: messageCollectionService,
		userProvider:             userProvider,
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
