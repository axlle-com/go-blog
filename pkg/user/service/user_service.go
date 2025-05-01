package service

import (
	"errors"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/user/models"
	"github.com/axlle-com/blog/pkg/user/repository"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

func NewUserService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
) *UserService {
	return &UserService{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

func (s *UserService) GetByEmail(email string) (user *models.User, err error) {
	user, err = s.userRepo.GetByEmail(email)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	if user == nil || user.ID == 0 {
		return nil, errors.New("invalid password or login")
	}

	return
}

func (s *UserService) Create(user *models.User) error {
	return s.userRepo.Create(user)
}

func (s *UserService) CreateFromInterface(user contracts.User) (*models.User, error) {
	if user == nil {
		return nil, nil
	}

	email, err := s.userRepo.GetByEmail(user.GetEmail())
	if errors.Is(err, gorm.ErrRecordNotFound) {
		uuid, err := s.userRepo.GetByUUID(user.GetUUID())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUser := &models.User{}
			newUser.FromInterface(user)
			err = s.userRepo.Create(newUser)
			if err != nil {
				return nil, err
			}
			return newUser, nil
		}

		if uuid == nil {
			return uuid, err
		}

		newUser := &models.User{}
		newUser.FromInterface(user)
		err = s.userRepo.Create(newUser)
		if err != nil {
			return nil, err
		}

		newUserHasUser, err := s.userRepo.GetRelation(uuid.UUID, newUser.UUID)
		if newUserHasUser != nil {
			return newUser, nil
		}

		newUserHasUser = &models.UserHasUser{}
		newUserHasUser.UserUUID = uuid.UUID
		newUserHasUser.RelationUUID = newUser.UUID
		err = s.userRepo.Attach(newUserHasUser)
		if err != nil {
			return newUser, err
		}
		return newUser, nil

	}

	if email == nil {
		return email, err
	}

	if email.UUID == user.GetUUID() {
		return email, nil
	}

	newUserHasUser, err := s.userRepo.GetRelation(email.UUID, user.GetUUID())
	if newUserHasUser != nil {
		return email, nil
	}

	newUserHasUser = &models.UserHasUser{}
	newUserHasUser.UserUUID = email.UUID
	newUserHasUser.RelationUUID = user.GetUUID()
	err = s.userRepo.Attach(newUserHasUser)
	if err != nil {
		return email, err
	}

	return email, err
}
