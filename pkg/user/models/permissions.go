package models

import "time"

type Permission struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"size:255;unique" json:"name"`
	CreatedAt *time.Time `gorm:"index" json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	Users []User `gorm:"many2many:user_has_permission;" json:"users,omitempty"`
	Roles []Role `gorm:"many2many:role_has_permission;" json:"roles,omitempty"`
}
