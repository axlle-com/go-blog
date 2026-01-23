package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UUID         uuid.UUID      `gorm:"type:uuid;index,using:hash;not null" json:"uuid" form:"uuid" binding:"-"`
	UserID       uint           `gorm:"index;not null" json:"user_id" form:"user_id" binding:"omitempty"`
	File         string         `gorm:"size:255;not null;unique" json:"file"`
	OriginalName string         `gorm:"size:255;not null" json:"original_name"`
	Size         int64          `gorm:"size:255;not null" json:"size,omitempty"`
	Type         string         `gorm:"size:255" json:"type,omitempty"`
	ReceivedAt   *time.Time     `gorm:"index" json:"received_at" form:"received_at" binding:"omitempty"`
	CreatedAt    *time.Time     `gorm:"index" json:"created_at,omitempty"`
	UpdatedAt    *time.Time     `json:"updated_at,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at" binding:"omitempty"`
}

func (f *File) GetTable() string {
	return "files"
}
func (f *File) Fields() []string {
	return []string{
		"id",
		"uuid",
		"user_id",
		"file",
		"original_name",
		"size",
		"type",
		"received_at",
	}
}

func (f *File) GetUUID() uuid.UUID {
	return f.UUID
}

func (f *File) SetUUID() {
	if f.UUID == uuid.Nil {
		f.UUID = uuid.New()
	}
}
