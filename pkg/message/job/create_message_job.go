package mailer

import (
	"context"
	"github.com/axlle-com/blog/app/models/contracts"
	contracts2 "github.com/axlle-com/blog/pkg/message/contracts"
	"github.com/axlle-com/blog/pkg/message/form"
	"time"
)

func NewCreateMessageJob(
	messageService contracts2.MessageService,
	form form.Form,
) contracts.Job {
	return &CreateMessageJob{
		messageService: messageService,
		form:           form,
		start:          time.Now(),
	}
}

type CreateMessageJob struct {
	messageService contracts2.MessageService
	form           form.Form
	start          time.Time
}

func (j *CreateMessageJob) Run(ctx context.Context) error {
	_, err := j.messageService.Create(j.form.Model(), j.form.GetUserUUID())
	if err != nil {
		return err
	}
	return nil
}

func (j *CreateMessageJob) GetData() []byte {
	return []byte(j.form.Data())
}

func (j *CreateMessageJob) GetName() string {
	return "CreateMessage"
}

func (j *CreateMessageJob) Duration() float64 {
	elapsed := time.Since(j.start)
	return float64(elapsed.Nanoseconds()) / 1e6
}
