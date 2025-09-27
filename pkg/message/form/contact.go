package form

import (
	"encoding/json"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/message/models"
)

type Contact struct {
	From     *string `json:"from" binding:"required,email"`
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
	if c.From != nil {
		return *c.From
	}
	return ""
}

func (c *Contact) GetUserUUID() string {
	return c.UserUUID
}

func (c *Contact) Title() string {
	return "Форма обратной связи"
}

func (c *Contact) Model() *models.Message {
	return &models.Message{
		From:    c.From,
		Subject: c.Subject,
		Body:    *c.Body,
	}
}
