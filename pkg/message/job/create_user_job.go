package mailer

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/message/form"
	"github.com/axlle-com/blog/pkg/message/models"
	"github.com/axlle-com/blog/pkg/user/provider"
)

func NewCreateUserJob(
	userProvider provider.UserProvider,
	form form.Form,
) contracts.Job {
	return &CreateUserJob{
		userProvider: userProvider,
		form:         form,
		start:        time.Now(),
	}
}

type CreateUserJob struct {
	userProvider provider.UserProvider
	form         form.Form
	start        time.Time
}

func (j *CreateUserJob) Run(ctx context.Context) error {
	parse, err := uuid.Parse(j.form.GetUserUUID())
	if err != nil {
		return err
	}

	user := &models.User{
		UUID:  parse,
		Email: j.form.GetFrom(),
	}
	_, err = j.userProvider.Create(user)
	if err != nil {
		return err
	}
	return nil
}

func (j *CreateUserJob) GetData() []byte {
	return []byte(j.form.Data())
}

func (j *CreateUserJob) GetName() string {
	return "CreateUser"
}

func (j *CreateUserJob) Duration() float64 {
	elapsed := time.Since(j.start)
	return float64(elapsed.Nanoseconds()) / 1e6
}
