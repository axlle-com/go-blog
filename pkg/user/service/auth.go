package service

import (
	"errors"
	"github.com/axlle-com/blog/pkg/common/config"
	. "github.com/axlle-com/blog/pkg/common/models"
	. "github.com/axlle-com/blog/pkg/user/http/models"
	"github.com/axlle-com/blog/pkg/user/repository"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

func Auth(authInput AuthInput) (userFound *User, err error) {
	userRepo := repository.NewRepository()
	userFound, err = userRepo.GetByEmail(authInput.Email)

	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	if userFound == nil || userFound.ID == 0 {
		return nil, errors.New("Invalid password or login")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.PasswordHash), []byte(authInput.Password)); err != nil {
		return nil, errors.New("Invalid password or login")
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(config.GetConfig().KeyJWT))
	if err != nil {
		return
	}
	userFound.AuthToken = &token

	return
}
