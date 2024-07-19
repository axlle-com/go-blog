package models

import "time"

type Template struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Title     string     `gorm:"size:255;not null" json:"title"`
	Name      string     `gorm:"size:45;not null" json:"name"`
	Resource  *string    `gorm:"size:255" json:"resource,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
