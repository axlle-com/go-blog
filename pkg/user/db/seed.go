package db

import (
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/bxcodec/faker/v3"
	"math/rand"
	"time"

	. "github.com/axlle-com/blog/pkg/user/models"
	. "github.com/axlle-com/blog/pkg/user/repository"
)

type seeder struct {
	user       UserRepository
	role       RoleRepository
	permission PermissionRepository
}

func NewSeeder(
	user UserRepository,
	role RoleRepository,
	permission PermissionRepository,
) contracts.Seeder {
	return &seeder{
		user:       user,
		role:       role,
		permission: permission,
	}
}

func (s *seeder) Seed() {
	s.seedPermissions()
	s.seedRoles()

	phone := "+7-900-111-22-33"
	createdAt := time.Now()
	updatedAt := time.Now()

	role, _ := s.role.GetByName("admin")
	user := User{
		Avatar:    db.StrPtr("/public/img/user.svg"),
		FirstName: "Admin",
		LastName:  "Admin",
		Phone:     &phone,
		Email:     "admin@admin.ru",
		IsEmail:   db.BoolToBoolPtr(true),
		IsPhone:   db.BoolToBoolPtr(true),
		Status:    10,
		Password:  "123456",
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
		Roles:     []Role{*role},
	}

	err := s.user.Create(&user)
	if err != nil {
		logger.Errorf("Failed to create user: %v", err.Error())
	}

	logger.Info("Database seeded User successfully!")
}

func (s *seeder) SeedTest(n int) {
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

		user := User{
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
			logger.Errorf("Failed to create user %d: %v", i, err.Error())
		}
	}
	logger.Info("Database seeded User successfully!")
}
