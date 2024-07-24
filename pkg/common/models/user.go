package models

import (
	"github.com/axlle-com/blog/pkg/common/config"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type User struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	FirstName          string     `gorm:"size:255;not null;default:'Undefined'" json:"first_name"`
	LastName           string     `gorm:"size:255;not null;default:'Undefined'" json:"last_name"`
	Patronymic         *string    `gorm:"size:255" json:"patronymic,omitempty"`
	Phone              *string    `gorm:"size:255;unique" json:"phone,omitempty"`
	Email              *string    `gorm:"size:255;unique" json:"email,omitempty"`
	IsEmail            *bool      `gorm:"default:false" json:"is_email,omitempty"`
	IsPhone            *bool      `gorm:"default:false" json:"is_phone,omitempty"`
	Status             int8       `gorm:"default:0" json:"status"`
	Avatar             *string    `gorm:"size:255" json:"avatar,omitempty"`
	Password           string     `gorm:"-" json:"password"`
	PasswordHash       string     `gorm:"size:255;not null" json:"password_hash"`
	RememberToken      *string    `gorm:"size:500" json:"remember_token,omitempty"`
	AuthToken          *string    `gorm:"size:500;default:null;index" json:"auth_token"`
	AuthKey            *string    `gorm:"size:32" json:"auth_key,omitempty"`
	PasswordResetToken *string    `gorm:"size:255;unique" json:"password_reset_token,omitempty"`
	CreatedAt          *time.Time `json:"created_at,omitempty"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
	DeletedAt          *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (u *User) SetPasswordHash() {
	if u.Password == "" {
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	u.PasswordHash = string(passwordHash)
}

func (u *User) SetAuthToken() (token string, err error) {
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  u.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err = generateToken.SignedString([]byte(config.GetConfig().KeyJWT))

	if err != nil {
		return
	}
	u.AuthToken = &token
	return
}
