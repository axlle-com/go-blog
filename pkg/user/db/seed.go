package db

import (
	"math/rand"
	"time"

	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/user/models"
	"github.com/axlle-com/blog/pkg/user/repository"
	"github.com/bxcodec/faker/v3"
)

type seeder struct {
	user       repository.UserRepository
	role       repository.RoleRepository
	permission repository.PermissionRepository
}

func NewSeeder(
	user repository.UserRepository,
	role repository.RoleRepository,
	permission repository.PermissionRepository,
) contract.Seeder {
	return &seeder{
		user:       user,
		role:       role,
		permission: permission,
	}
}

func (s *seeder) Seed() error {
	err := s.seedPermissions()
	if err != nil {
		return err
	}

	err = s.seedRoles()
	if err != nil {
		return err
	}

	email, _ := s.user.GetByEmail("admin@ax-box.ru")
	if email != nil {
		logger.Info("[Useer][seeder][Seed] User is full")
		return nil
	}

	phone := "+7-900-111-22-33"
	createdAt := time.Now()
	updatedAt := time.Now()

	role, _ := s.role.GetByName("superadmin")
	user := models.User{
		Avatar:    db.StrPtr("/public/img/user.svg"),
		FirstName: "Admin",
		LastName:  "Admin",
		Phone:     &phone,
		Email:     "admin@ax-box.ru",
		IsEmail:   db.BoolToBoolPtr(true),
		IsPhone:   db.BoolToBoolPtr(true),
		Status:    10,
		Password:  "123456",
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
		Roles:     []models.Role{*role},
	}

	err = s.user.Create(&user)
	if err != nil {
		return err
	}

	logger.Info("[Useer][seeder][Seed] Database seeded User successfully!")
	return err
}

func (s *seeder) SeedTest(n int) error {
	for i := 0; i < n; i++ {
		firstName := faker.FirstName()
		lastName := faker.LastName()
		phone := faker.Phonenumber()
		password := faker.Password()
		rememberToken := faker.UUIDDigit()
		authKey := faker.UUIDDigit()
		authToken := faker.UUIDHyphenated()
		passwordResetToken := faker.UUIDHyphenated()
		createdAt := time.Now()
		updatedAt := time.Now()

		user := models.User{
			Avatar:             db.StrPtr("/public/img/user.svg"),
			FirstName:          firstName,
			LastName:           lastName,
			Phone:              &phone,
			Email:              faker.Email(),
			IsEmail:            db.IntToBoolPtr(),
			IsPhone:            db.IntToBoolPtr(),
			Status:             int8(rand.Intn(10)),
			Password:           password,
			RememberToken:      &rememberToken,
			AuthKey:            &authKey,
			AuthToken:          &authToken,
			PasswordResetToken: &passwordResetToken,
			CreatedAt:          &createdAt,
			UpdatedAt:          &updatedAt,
		}
		err := s.user.Create(&user)
		if err != nil {
			logger.Errorf("[Useer][seeder][SeedTest] Failed to create user %d: %v", i, err.Error())
		}
	}
	logger.Info("[Useer][seeder][SeedTest] Database seeded User successfully!")
	return nil
}
