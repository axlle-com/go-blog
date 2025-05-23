package mailer

import (
	"context"
	config2 "github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/models/contracts"
	"gopkg.in/gomail.v2"
	"time"
)

func NewMailerJob(
	message contracts.MailRequest,
) contracts.Job {
	return &MailerJob{
		message: message,
		start:   time.Now(),
	}
}

type MailerJob struct {
	message contracts.MailRequest
	start   time.Time
}

func (j *MailerJob) Run(ctx context.Context) error {
	config := config2.Config()
	if !config.SMTPActive() {
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", j.message.From())
	m.SetHeader("To", j.message.To())
	m.SetHeader("Subject", j.message.Subject())
	m.SetBody("text/html", j.message.Body())

	d := gomail.NewDialer(config.SMTPHost(), config.SMTPPort(), config.SMTPUsername(), config.SMTPPassword())

	err := d.DialAndSend(m)
	if err != nil {
		return err
	}

	return nil
}

func (j *MailerJob) GetData() []byte {
	return []byte(j.message.ToString())
}

func (j *MailerJob) GetName() string {
	return "Mailer"
}

func (j *MailerJob) Duration() float64 {
	elapsed := time.Since(j.start)
	return float64(elapsed.Nanoseconds()) / 1e6
}
