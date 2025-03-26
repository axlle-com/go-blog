package service

import (
	"errors"
	"gorm.io/gorm"

	. "github.com/axlle-com/blog/pkg/user/models"
	. "github.com/axlle-com/blog/pkg/user/repository"
)

type UserService struct {
	userRepo UserRepository
}

func NewUserService(
	userRepo UserRepository,
) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetByEmail(email string) (user *User, err error) {
	user, err = s.userRepo.GetByEmail(email)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	if user == nil || user.ID == 0 {
		return nil, errors.New("invalid password or login")
	}

	return
}

func (s *UserService) Create(user *User) error {
	return s.userRepo.Create(user)
}
