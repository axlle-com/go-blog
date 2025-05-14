package models

import "github.com/google/uuid"

type PostTagHasResource struct {
	PostTagID    uint      `gorm:"index;not null"`
	ResourceUUID uuid.UUID `gorm:"type:uuid;index,using:hash" json:"resource_uuid" form:"resource_uuid" binding:"-"`
	Sort         int       `gorm:"index;not null"`
}
