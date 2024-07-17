package models

import (
	"github.com/axlle-com/blog/pkg/common/config"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"time"
)

type User struct {
	ID                 uint           `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	FirstName          string         `gorm:"type:varchar(255);not null;default:'Undefined';index" json:"first_name"`
	LastName           string         `gorm:"type:varchar(255);not null;default:'Undefined';index" json:"last_name"`
	Patronymic         string         `gorm:"type:varchar(255);default:null" json:"patronymic"`
	Phone              string         `gorm:"type:varchar(255);default:null;unique" json:"phone"`
	Email              string         `gorm:"type:varchar(255);default:null;unique" json:"email"`
	IsEmail            uint8          `gorm:"type:smallint;default:0" json:"is_email"`
	IsPhone            uint8          `gorm:"type:smallint;default:0" json:"is_phone"`
	Status             int16          `gorm:"type:smallint;not null;default:0" json:"status"`
	Avatar             string         `gorm:"type:varchar(255);default:null" json:"avatar"`
	Password           string         `gorm:"-" json:"password"`
	PasswordHash       string         `gorm:"type:varchar(255);not null" json:"password_hash"`
	RememberToken      string         `gorm:"type:varchar(500);default:null" json:"remember_token"`
	AuthKey            string         `gorm:"type:varchar(32);default:null" json:"auth_key"`
	AuthToken          string         `gorm:"type:varchar(32);default:null;index" json:"auth_token"`
	PasswordResetToken string         `gorm:"type:varchar(255);default:null;unique" json:"password_reset_token"`
	CreatedAt          time.Time      `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index;type:timestamp;default:null;index" json:"deleted_at"`
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
	u.AuthToken = token
	return
}
