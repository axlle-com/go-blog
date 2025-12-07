package models

import (
	"fmt"
	"time"

	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/models/contract"
)

type Template struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	UserID       *uint      `gorm:"index" json:"user_id" form:"user_id" binding:"omitempty"`
	Title        string     `gorm:"size:255;not null" json:"title"`
	IsMain       bool       `gorm:"index;not null;default:false" json:"is_main" form:"is_main" binding:"omitempty"`
	Theme        *string    `gorm:"size:255;not null" json:"theme"`
	Name         string     `gorm:"size:255;not null" json:"name"`
	ResourceName *string    `gorm:"size:255" json:"resource_name,omitempty"`
	HTML         *string    `gorm:"type:text" json:"html" binding:"omitempty"`
	JS           *string    `gorm:"type:text" json:"js,omitempty"`
	CSS          *string    `gorm:"type:text" json:"css,omitempty"`
	CreatedAt    *time.Time `gorm:"index" json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	User contract.User `gorm:"-" json:"user"`
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
	if t.Name == "" {
		return "default"
	}
	return t.Name
}

func (t *Template) GetFullName(resourceName string) string {
	return fmt.Sprintf("%s.%s", resourceName, t.GetName())
}

func (t *Template) GetResourceName() string {
	if t.ResourceName != nil {
		return *t.ResourceName
	}

	return ""
}

func (t *Template) GetThemeName() string {
	if t.Theme != nil {
		return *t.Theme
	}

	return "default"
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

func (t *Template) Saving() {
	if t.Theme == nil || *t.Theme == "" {
		v := config.Config().Layout()
		t.Theme = &v
	}
}

func (t *Template) Creating() {
	t.Saving()
}

func (t *Template) Updating() {
	t.Saving()
}
