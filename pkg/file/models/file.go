package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type File struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UUID         uuid.UUID      `gorm:"type:uuid;index,using:hash;not null" json:"uuid" form:"uuid" binding:"-"`
	UserID       uint           `gorm:"index;not null" json:"user_id" form:"user_id" binding:"omitempty"`
	File         string         `gorm:"size:255;not null;unique" json:"file"`
	OriginalName string         `gorm:"size:255;not null" json:"original_name"`
	Size         int64          `gorm:"size:255;not null" json:"size,omitempty"`
	Type         string         `gorm:"size:255" json:"type,omitempty"`
	IsReceived   bool           `gorm:"index;not null;default:false" json:"is_received" form:"is_received" binding:"omitempty"`
	CreatedAt    *time.Time     `gorm:"index" json:"created_at,omitempty"`
	UpdatedAt    *time.Time     `json:"updated_at,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at" binding:"omitempty"`
}

func (u *File) Fields() []string {
	return []string{
		"id",
		"uuid",
		"user_id",
		"file",
		"original_name",
		"size",
		"type",
		"is_received",
	}
}

func (u *File) GetUUID() uuid.UUID {
	return u.UUID
}

func (u *File) SetUUID() {
	if u.UUID == uuid.Nil {
		u.UUID = uuid.New()
	}
}
