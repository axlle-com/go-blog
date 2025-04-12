package models

import (
	"fmt"
	"github.com/axlle-com/blog/app/models/contracts"
	"time"
)

type Template struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    *uint      `json:"user_id" form:"user_id" binding:"omitempty"`
	Title     string     `gorm:"size:255;not null" json:"title"`
	Name      string     `gorm:"size:45;not null" json:"name"`
	Tabular   *string    `gorm:"size:255" json:"tabular,omitempty"`
	JS        *string    `gorm:"type:text" json:"js,omitempty"`
	CSS       *string    `gorm:"type:text" json:"css,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	User contracts.User `gorm:"-" json:"user"`
}

func (t *Template) GetTable() string {
	return "templates"
}

func (t *Template) AdminURL() string {
	if t.ID == 0 {
		return "/admin/templates"
	}
	return fmt.Sprintf("/admin/templates/%d", t.ID)
}

func (t *Template) GetID() uint {
	return t.ID
}

func (t *Template) GetTitle() string {
	return t.Title
}

func (t *Template) GetName() string {
	return t.Title
}

func (t *Template) GetTabular() string {
	return *t.Tabular
}
