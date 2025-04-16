package models

import (
	"fmt"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserUUID   uuid.UUID  `gorm:"type:uuid;index,using:hash" json:"user_uuid" form:"user_uuid" binding:"-"`
	From       *string    `gorm:"size:255null" json:"from" form:"from" binding:"omitempty"`
	To         *string    `gorm:"size:255;null" json:"to" form:"to" binding:"omitempty"`
	Subject    *string    `gorm:"size:255;not null" json:"subject" form:"subject" binding:"required,max=255"`
	Body       string     `gorm:"type:text" json:"body" form:"body" binding:"omitempty"`
	Attachment string     `gorm:"type:text" json:"attachment" form:"attachment" binding:"omitempty"`
	Viewed     bool       `gorm:"default:false" json:"viewed"  form:"viewed" binding:"omitempty"`
	CreatedAt  *time.Time `gorm:"index" json:"created_at" form:"created_at" binding:"omitempty"`
	UpdatedAt  *time.Time `json:"updated_at" form:"updated_at" binding:"omitempty"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at" form:"deleted_at" binding:"omitempty"`

	User contracts.User `gorm:"-" json:"user"`
}

func (m *Message) GetTable() string {
	return "messages"
}

func (m *Message) AdminURL() string {
	if m.ID == 0 {
		return "/admin/messages"
	}
	return fmt.Sprintf("/admin/messages/%d", m.ID)
}

func (m *Message) GetID() uint {
	return m.ID
}

func (m *Message) UserLastName() string {
	var lastName string
	if m.User != nil {
		lastName = m.User.GetLastName()
	}
	return lastName
}

func (m *Message) Date() string {
	if m.CreatedAt == nil {
		return ""
	}
	return m.CreatedAt.Format("02.01.2006 15:04:05")
}
