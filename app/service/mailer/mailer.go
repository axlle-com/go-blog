package mailer

import (
	"github.com/axlle-com/blog/app/models/contracts"
	"gopkg.in/gomail.v2"
)

type Smtp struct {
	config contracts.Config
	queue  contracts.Queue
}

func NewMailer(config contracts.Config, queue contracts.Queue) contracts.Mailer {
	return &Smtp{
		config: config,
		queue:  queue,
	}
}

func (s *Smtp) SendMail(message contracts.MailRequest) error {
	if !s.config.SMTPActive() {
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.config.SMTPUsername())
	m.SetHeader("To", message.GetTo())
	m.SetHeader("Subject", message.GetSubject())
	m.SetBody("text/html", message.GetBody())

	d := gomail.NewDialer(s.config.SMTPHost(), s.config.SMTPPort(), s.config.SMTPUsername(), s.config.SMTPPassword())

	err := d.DialAndSend(m)
	if err != nil {
		return err
	}

	return nil
}
