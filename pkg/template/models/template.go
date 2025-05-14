package models

import (
	"fmt"
	"github.com/axlle-com/blog/app/models/contracts"
	"time"
)

type Template struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	UserID       *uint      `gorm:"index" json:"user_id" form:"user_id" binding:"omitempty"`
	Title        string     `gorm:"size:255;not null" json:"title"`
	IsMain       bool       `gorm:"index;not null;default:false" json:"is_main" form:"is_main" binding:"omitempty"`
	Name         string     `gorm:"size:45;not null;unique" json:"name"`
	ResourceName *string    `gorm:"size:255" json:"resource_name,omitempty"`
	HTML         *string    `gorm:"type:text" json:"html" binding:"omitempty"`
	JS           *string    `gorm:"type:text" json:"js,omitempty"`
	CSS          *string    `gorm:"type:text" json:"css,omitempty"`
	CreatedAt    *time.Time `gorm:"index" json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`

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
	return t.Name
}

func (t *Template) GetResourceName() string {
	return *t.ResourceName
}

func (t *Template) UserLastName() string {
	var lastName string
	if t.User != nil {
		lastName = t.User.GetLastName()
	}
	return lastName
}

func (t *Template) Date() string {
	if t.CreatedAt == nil {
		return ""
	}
	return t.CreatedAt.Format("02.01.2006 15:04:05")
}
