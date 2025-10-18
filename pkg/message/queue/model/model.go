package model

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/message/models"
	"github.com/google/uuid"
)

type Message struct {
	UserUUID   string  `json:"user_uuid" form:"user_uuid" binding:"-"`
	From       *string `json:"from" form:"from" binding:"omitempty"`
	To         *string `json:"to" form:"to" binding:"omitempty"`
	Subject    *string `json:"subject" form:"subject" binding:"required,max=255"`
	Body       string  `json:"body" form:"body" binding:"omitempty"`
	Attachment string  `json:"attachment" form:"attachment" binding:"omitempty"`
}

func (m *Message) Model() *models.Message {
	var userUUID uuid.UUID
	if m.UserUUID != "" {
		newUUID, err := uuid.Parse(m.UserUUID)
		if err != nil {
			logger.Errorf("[message][Model] Invalid UUID: %v", err)
		}
		userUUID = newUUID
	}
	return &models.Message{
		From:     m.From,
		Subject:  m.Subject,
		Body:     m.Body,
		UserUUID: userUUID,
	}
}
