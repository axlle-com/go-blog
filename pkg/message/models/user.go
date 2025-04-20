package models

import (
	"github.com/google/uuid"
	"regexp"
	"strings"
)

type User struct {
	UUID       uuid.UUID `gorm:"type:uuid;index,using:hash" json:"uuid" form:"uuid" binding:"-"`
	FirstName  string    `gorm:"size:255;not null;default:'Undefined'" json:"first_name"`
	LastName   string    `gorm:"size:255;not null;default:'Undefined'" json:"last_name"`
	Patronymic *string   `gorm:"size:255" json:"patronymic,omitempty"`
	Phone      *string   `gorm:"size:255;unique" json:"phone,omitempty"`
	Email      string    `gorm:"size:255;unique;not null" json:"email"`
}

func (u *User) SetPhone() {
	if u.Phone != nil {
		if ok, phone := ValidateAndCleanPhone(*u.Phone); ok {
			u.Phone = &phone
			return
		}
	}
	u.Phone = nil
}

func ValidateAndCleanPhone(phone string) (bool, string) {
	re := regexp.MustCompile(`\D`)
	cleanedPhone := re.ReplaceAllString(phone, "")

	if len(cleanedPhone) != 11 {
		if len(cleanedPhone) == 10 {
			cleanedPhone = "7" + cleanedPhone
			return true, cleanedPhone
		}
		return false, cleanedPhone
	}

	if cleanedPhone[0] != '7' && cleanedPhone[0] != '8' {
		return false, cleanedPhone
	}

	if cleanedPhone[0] == '8' {
		cleanedPhone = "7" + cleanedPhone[1:]
	}

	return true, cleanedPhone
}

func (u *User) GetID() uint {
	return 0
}

func (u *User) GetFirstName() string {
	return u.FirstName
}

func (u *User) GetLastName() string {
	return u.LastName
}

func (u *User) GetPatronymic() string {
	if u.Patronymic != nil {
		return *u.Patronymic
	}
	return ""
}

func (u *User) GetFullName() string {
	var parts []string

	if s := strings.TrimSpace(u.FirstName); s != "" {
		parts = append(parts, s)
	}
	if s := strings.TrimSpace(u.LastName); s != "" {
		parts = append(parts, s)
	}
	if u.Patronymic != nil {
		if s := strings.TrimSpace(*u.Patronymic); s != "" {
			parts = append(parts, s)
		}
	}

	return strings.Join(parts, " ")
}

func (u *User) GetPhone() string {
	if u.Phone != nil {
		return *u.Phone
	}
	return ""
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetStatus() int8 {
	return 0
}

func (u *User) GetRoles() []string {
	return nil
}

func (u *User) GetPermissions() []string {
	return nil
}

func (u *User) GetUUID() uuid.UUID {
	return u.UUID
}
