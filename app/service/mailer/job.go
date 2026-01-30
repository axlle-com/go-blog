package mailer

import (
	"context"
	"time"

	"github.com/axlle-com/blog/app/models/contract"
	"gopkg.in/gomail.v2"
)

func NewMailerJob(
	config contract.Config,
	message contract.MailRequest,
) contract.Job {
	return &MailerJob{
		config:  config,
		message: message,
		start:   time.Now(),
	}
}

type MailerJob struct {
	config  contract.Config
	message contract.MailRequest
	start   time.Time
}

func (j *MailerJob) Run(ctx context.Context) error {
	if !j.config.SMTPActive() {
		return nil
	}

	message := gomail.NewMessage()
	message.SetHeader("From", j.message.GetFrom())
	message.SetHeader("To", j.message.GetTo())
	message.SetHeader("Subject", j.message.GetSubject())
	message.SetBody("text/html", j.message.GetBody())

	dialer := gomail.NewDialer(
		j.config.SMTPHost(),
		j.config.SMTPPort(),
		j.config.SMTPUsername(),
		j.config.SMTPPassword(),
	)

	return dialer.DialAndSend(message)
}

func (j *MailerJob) GetData() []byte {
	return []byte(j.message.ToString())
}

func (j *MailerJob) GetName() string {
	return "mailer"
}

func (j *MailerJob) GetQueue() string {
	return "mailer"
}

func (j *MailerJob) GetAction() string {
	return "send"
}

func (j *MailerJob) Duration() float64 {
	elapsed := time.Since(j.start)

	return float64(elapsed.Nanoseconds()) / 1e6
}
