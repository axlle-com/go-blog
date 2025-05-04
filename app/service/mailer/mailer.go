package mailer

import (
	"github.com/axlle-com/blog/app/models/contracts"
)

type Smtp struct {
	queue contracts.Queue
}

func NewMailer(queue contracts.Queue) contracts.Mailer {
	return &Smtp{
		queue: queue,
	}
}

func (s *Smtp) SendMail(message contracts.MailRequest) {
	s.queue.Enqueue(NewMailerJob(message), 0)
}
