package mailer

import (
	"context"
	"encoding/json"
	"time"

	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/message/form"
)

func NewCreateUserJob(
	form form.Form,
) contracts.Job {
	return &CreateUserJob{
		form:  form,
		start: time.Now(),
	}
}

type CreateUserJob struct {
	form  form.Form
	start time.Time
}

func (j *CreateUserJob) Run(ctx context.Context) error {
	return nil
}

func (j *CreateUserJob) GetData() []byte {
	payload := struct {
		UUID  string `json:"uuid"`
		Email string `json:"email"`
	}{
		j.form.GetUserUUID(),
		j.form.GetFrom(),
	}

	bytes, _ := json.Marshal(payload)
	return app.NewEnvelopeQueue().ConvertData("create", string(bytes))
}

func (j *CreateUserJob) GetName() string {
	return "users"
}

func (j *CreateUserJob) GetQueue() string {
	return "users"
}

func (j *CreateUserJob) GetAction() string {
	return "create"
}

func (j *CreateUserJob) Duration() float64 {
	elapsed := time.Since(j.start)
	return float64(elapsed.Nanoseconds()) / 1e6
}
