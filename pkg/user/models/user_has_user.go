package models

import (
	"github.com/google/uuid"
	"time"
)

type UserHasUser struct {
	UserUUID     uuid.UUID `gorm:"type:uuid;index,using:hash" json:"user_uuid" form:"user_uuid" binding:"-"`
	RelationUUID uuid.UUID `gorm:"type:uuid;index,using:hash" json:"relation_uuid" form:"relation_uuid" binding:"-"`

	CreatedAt *time.Time `gorm:"index" json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
