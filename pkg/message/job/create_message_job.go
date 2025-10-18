package mailer

import (
	"context"
	"time"

	"github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/message/form"
)

func NewCreateMessageJob(
	form form.Form,
) contracts.Job {
	return &CreateMessageJob{
		form:  form,
		start: time.Now(),
	}
}

type CreateMessageJob struct {
	form  form.Form
	start time.Time
}

func (j *CreateMessageJob) Run(ctx context.Context) error {
	return nil
}

func (j *CreateMessageJob) GetData() []byte {
	return models.NewEnvelopeQueue().ConvertData("create", j.form.Data())
}

func (j *CreateMessageJob) GetName() string {
	return "messages"
}

func (j *CreateMessageJob) GetQueue() string {
	return "messages"
}

func (j *CreateMessageJob) GetAction() string {
	return "create"
}

func (j *CreateMessageJob) Duration() float64 {
	elapsed := time.Since(j.start)
	return float64(elapsed.Nanoseconds()) / 1e6
}
