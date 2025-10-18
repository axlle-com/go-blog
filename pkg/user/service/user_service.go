package service

import (
	"errors"

	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/user/models"
	"github.com/axlle-com/blog/pkg/user/repository"
	"github.com/google/uuid"
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

	byEmail, err := s.userRepo.GetByEmail(user.GetEmail())
	if errors.Is(err, gorm.ErrRecordNotFound) {
		byUUID, err := s.userRepo.GetByUUID(user.GetUUID())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUser := &models.User{}
			newUser.FromInterface(user)
			err = s.userRepo.Create(newUser)
			if err != nil {
				return nil, err
			}
			return newUser, nil
		}

		if byUUID == nil {
			return byUUID, err
		}

		newUser := &models.User{}
		newUser.FromInterface(user)
		newUser.UUID = uuid.New()
		err = s.userRepo.Create(newUser)
		if err != nil {
			return nil, err
		}

		newUserHasUser, err := s.userRepo.GetRelation(byUUID.UUID, newUser.UUID)
		if newUserHasUser != nil {
			return newUser, nil
		}

		newUserHasUser = &models.UserHasUser{}
		newUserHasUser.UserUUID = byUUID.UUID
		newUserHasUser.RelationUUID = newUser.UUID
		err = s.userRepo.Attach(newUserHasUser)
		if err != nil {
			return newUser, err
		}
		return newUser, nil

	}

	if byEmail == nil {
		return byEmail, err
	}

	if byEmail.UUID == user.GetUUID() {
		return byEmail, nil
	}

	newUserHasUser, err := s.userRepo.GetRelation(byEmail.UUID, user.GetUUID())
	if newUserHasUser != nil {
		return byEmail, nil
	} else if err != nil {
		return byEmail, err
	}

	newUserHasUser = &models.UserHasUser{}
	newUserHasUser.UserUUID = byEmail.UUID
	newUserHasUser.RelationUUID = user.GetUUID()
	err = s.userRepo.Attach(newUserHasUser)
	if err != nil {
		return byEmail, err
	}

	return byEmail, err
}
