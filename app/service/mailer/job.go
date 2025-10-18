package mailer

import (
	"context"
	"time"

	"gopkg.in/gomail.v2"

	appConfig "github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/models/contracts"
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
	config := appConfig.Config()
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
