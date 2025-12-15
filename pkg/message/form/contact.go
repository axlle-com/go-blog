package form

import (
	"encoding/json"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/message/models"
)

type Contact struct {
	Email    *string `json:"email" binding:"required,email"`
	UserName *string `json:"user_name" binding:"required"`
	To       *string `json:"to" binding:"omitempty,email"`
	Subject  *string `json:"subject" binding:"required"`
	Body     *string `json:"body" binding:"required"`
	UserUUID string  `json:"user_uuid" binding:"omitempty"`
}

func (c *Contact) Data() string {
	data, err := json.Marshal(*c)
	if err != nil {
		logger.Errorf("[Contact][Data]Error: %v", err)
	}
	return string(data)
}

func (c *Contact) Name() string {
	return "contact"
}

func (c *Contact) GetFrom() string {
	if c.Email != nil {
		return *c.Email
	}
	return ""
}

func (c *Contact) GetTo() string {
	if c.To != nil {
		return *c.To
	}
	return ""
}

func (c *Contact) GetUserUUID() string {
	return c.UserUUID
}

func (c *Contact) GetUserName() *string {
	return c.UserName
}

func (c *Contact) Title() string {
	return "Форма обратной связи"
}

func (c *Contact) Model() *models.Message {
	return &models.Message{
		From:    c.Email,
		Subject: c.Subject,
		Body:    *c.Body,
	}
}
