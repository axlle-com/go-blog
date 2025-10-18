package model

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/user/models"
	"github.com/google/uuid"
)

type User struct {
	UUID  string `json:"uuid" form:"uuid"`
	Email string `json:"email" form:"uuid"`
}

func (m *User) Model() *models.User {
	var userUUID uuid.UUID
	if m.UUID != "" {
		newUUID, err := uuid.Parse(m.UUID)
		if err != nil {
			logger.Errorf("[user][Model] Invalid UUID: %v", err)
		}
		userUUID = newUUID
	}
	return &models.User{
		Email: m.Email,
		UUID:  userUUID,
	}
}
