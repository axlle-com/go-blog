package models

import (
	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"regexp"
	"strings"
	"time"
)

type User struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	UUID               uuid.UUID      `gorm:"type:uuid;index,using:hash" json:"uuid" form:"uuid" binding:"-"`
	FirstName          string         `gorm:"size:255;not null;default:'Undefined'" json:"first_name"`
	LastName           string         `gorm:"size:255;not null;default:'Undefined'" json:"last_name"`
	Patronymic         *string        `gorm:"size:255" json:"patronymic,omitempty"`
	Phone              *string        `gorm:"size:255;unique" json:"phone,omitempty"`
	Email              string         `gorm:"size:255;unique;not null" json:"email"`
	IsEmail            *bool          `gorm:"default:false" json:"is_email,omitempty"`
	IsPhone            *bool          `gorm:"default:false" json:"is_phone,omitempty"`
	Status             int8           `gorm:"index;not null;default:0" json:"status"`
	Avatar             *string        `gorm:"size:255" json:"avatar,omitempty"`
	Password           string         `gorm:"-" json:"-"`
	PasswordHash       *string        `gorm:"size:255" json:"-"`
	RememberToken      *string        `gorm:"size:500;default:null;index" json:"-"`
	AuthToken          *string        `gorm:"size:500;default:null;index" json:"-"`
	AuthKey            *string        `gorm:"size:32;default:null;" json:"-"`
	PasswordResetToken *string        `gorm:"size:255;unique" json:"-"`
	CreatedAt          *time.Time     `gorm:"index" json:"created_at,omitempty"`
	UpdatedAt          *time.Time     `json:"updated_at,omitempty"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at" binding:"omitempty"`

	Roles       []Role       `gorm:"many2many:user_has_role;" json:"roles,omitempty"`
	Permissions []Permission `gorm:"many2many:user_has_permission;" json:"permissions,omitempty"`
}

func (u *User) Fields() []string {
	return []string{
		"id",
		"uuid",
		"first_name",
		"last_name",
		"patronymic",
		"phone",
		"email",
		"is_email",
		"is_phone",
		"status",
		"avatar",
		"created_at",
		"updated_at",
	}
}

func (u *User) SetUUID() {
	if u.UUID == uuid.Nil {
		u.UUID = uuid.New()
	}
}

func (u *User) Creating() {
	u.SetPasswordHash()
	u.SetPhone()
	u.SetUUID()
}

func (u *User) Updating() {
	u.SetPasswordHash()
	u.SetPhone()
}

func (u *User) SetPasswordHash() {
	if u.Password == "" {
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err)
	}

	str := string(passwordHash)
	if str == "" {
		u.PasswordHash = nil
	}
	u.PasswordHash = &str
}

func (u *User) SetAuthToken() (token string, err error) {
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  u.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err = generateToken.SignedString(config.Config().KeyJWT())

	if err != nil {
		logger.Error(err)
		return
	}
	u.AuthToken = &token
	return
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
	return u.ID
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
	return u.Status
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

func (u *User) FromInterface(user contracts.User) {
	u.FirstName = user.GetFirstName()
	u.LastName = user.GetLastName()
	u.Email = user.GetEmail()
	u.UUID = user.GetUUID()

	patronymic := user.GetPatronymic()
	if patronymic == "" {
		u.Patronymic = nil
	} else {
		u.Patronymic = &patronymic
	}

	phone := user.GetPhone()
	if phone == "" {
		u.Phone = nil
	} else {
		u.Phone = &phone
	}
}
