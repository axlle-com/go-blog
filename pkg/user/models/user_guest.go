package models

import (
	"github.com/google/uuid"
	"time"
)

type UserGuest struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UUID       uuid.UUID  `gorm:"type:uuid;index,using:hash" json:"uuid" form:"uuid" binding:"-"`
	FirstName  string     `gorm:"size:255;not null;default:'Undefined'" json:"first_name"`
	LastName   string     `gorm:"size:255;not null;default:'Undefined'" json:"last_name"`
	Patronymic *string    `gorm:"size:255" json:"patronymic,omitempty"`
	Phone      *string    `gorm:"size:255;unique" json:"phone,omitempty"`
	Email      string     `gorm:"size:255;unique;not null" json:"email"`
	CreatedAt  *time.Time `gorm:"index" json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	DeletedAt  *time.Time `gorm:"index" json:"-"`
}

func (u *UserGuest) Fields() []string {
	return []string{
		"id",
		"uuid",
		"first_name",
		"last_name",
		"patronymic",
		"phone",
		"email",
		"created_at",
		"updated_at",
	}
}

func (u *UserGuest) SetUUID() {
	if u.UUID == uuid.Nil {
		u.UUID = uuid.New()
	}
}

func (u *UserGuest) Creating() {
	u.SetPhone()
	u.SetUUID()
}

func (u *UserGuest) Updating() {
	u.SetPhone()
}

func (u *UserGuest) SetPhone() {
	if u.Phone != nil {
		if ok, phone := ValidateAndCleanPhone(*u.Phone); ok {
			u.Phone = &phone
			return
		}
	}
	u.Phone = nil
}

func (u *UserGuest) GetID() uint {
	return u.ID
}

func (u *UserGuest) GetFirstName() string {
	return u.FirstName
}

func (u *UserGuest) GetLastName() string {
	return u.LastName
}

func (u *UserGuest) GetPatronymic() string {
	return *u.Patronymic
}

func (u *UserGuest) GetPhone() string {
	return *u.Phone
}

func (u *UserGuest) GetEmail() string {
	return u.Email
}

func (u *UserGuest) GetUUID() uuid.UUID {
	return u.UUID
}

func (u *UserGuest) GetStatus() int8 {
	return 0
}

func (u *UserGuest) GetRoles() []string {
	return nil
}

func (u *UserGuest) GetPermissions() []string {
	return nil
}
