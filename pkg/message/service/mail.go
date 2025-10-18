package service

import (
	"github.com/axlle-com/blog/app/models/contracts"
	selfContracts "github.com/axlle-com/blog/pkg/message/contracts"
	"github.com/axlle-com/blog/pkg/message/form"
	mailer "github.com/axlle-com/blog/pkg/message/job"
	"github.com/axlle-com/blog/pkg/user/provider"
)

type MailService struct {
	messageService           selfContracts.MessageService
	messageCollectionService selfContracts.MessageCollectionService
	userProvider             provider.UserProvider
	mailer                   contracts.Mailer
	queue                    contracts.Queue
}

func NewMailService(
	messageService selfContracts.MessageService,
	messageCollectionService selfContracts.MessageCollectionService,
	userProvider provider.UserProvider,
	mailer contracts.Mailer,
	queue contracts.Queue,
) *MailService {
	return &MailService{
		messageService:           messageService,
		messageCollectionService: messageCollectionService,
		userProvider:             userProvider,
		mailer:                   mailer,
		queue:                    queue,
	}
}

func (s *MailService) SendContact(form form.Form) {
	mailRequest := NewInformer(form.Data(), form.Title())
	messageJob := mailer.NewCreateMessageJob(form)
	userJob := mailer.NewCreateUserJob(form)

	s.mailer.SendMail(mailRequest)
	s.queue.Enqueue(messageJob, 0)
	s.queue.Enqueue(userJob, 0)
}
