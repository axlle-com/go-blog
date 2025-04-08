package models

import "github.com/google/uuid"

type InfoBlockHasResource struct {
	ID           uint      `gorm:"primaryKey" json:"id" form:"id" binding:"-"`
	InfoBlockID  uint      `gorm:"index;not null"`
	ResourceUUID uuid.UUID `gorm:"type:uuid;index,using:hash" json:"uuid" form:"uuid" binding:"-"`
	Sort         int       `gorm:"index;not null"`
	Position     string    `gorm:"index"`
}
