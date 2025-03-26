package service

import (
	"errors"
	"github.com/axlle-com/blog/pkg/app/config"
	. "github.com/axlle-com/blog/pkg/user/http/models"
	user "github.com/axlle-com/blog/pkg/user/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	userService *UserService
}

func NewAuthService(
	userService *UserService,
) *AuthService {
	return &AuthService{userService: userService}
}

func (s *AuthService) Auth(authInput AuthInput) (user *user.User, err error) {
	user, err = s.userService.GetByEmail(authInput.Email)

	if err != nil {
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(authInput.Password)); err != nil {
		return nil, errors.New("invalid password or login")
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString(config.Config().KeyJWT())
	if err != nil {
		return
	}
	user.AuthToken = &token

	return
}
