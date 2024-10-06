package db

import (
	. "github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/user/models"
	. "github.com/axlle-com/blog/pkg/user/repository"
	"github.com/bxcodec/faker/v3"
	"log"
	"math/rand"
	"time"
)

func SeedUsers(n int) {
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
			Avatar:             StrPtr("/public/img/user.svg"),
			FirstName:          firstName,
			LastName:           lastName,
			Phone:              &phone,
			Email:              faker.Email(),
			IsEmail:            IntToBoolPtr(),
			IsPhone:            IntToBoolPtr(),
			Status:             int8(rand.Intn(10)),
			Password:           password,
			RememberToken:      &rememberToken,
			AuthKey:            &authKey,
			AuthToken:          &authToken,
			PasswordResetToken: &passwordResetToken,
			CreatedAt:          &createdAt,
			UpdatedAt:          &updatedAt,
		}
		err := NewRepo().Create(&user)
		if err != nil {
			log.Printf("Failed to create user %d: %v", i, err.Error())
		}
	}

	phone := "+7-900-111-22-33"
	rememberToken := faker.UUIDDigit()
	authKey := faker.UUIDDigit()
	authToken := faker.UUIDHyphenated()
	passwordResetToken := faker.UUIDHyphenated()
	createdAt := time.Now()
	updatedAt := time.Now()

	role, _ := NewRoleRepository().GetByName("admin")
	user := models.User{
		Avatar:             StrPtr("/public/img/user.svg"),
		FirstName:          "Admin",
		LastName:           "Admin",
		Phone:              &phone,
		Email:              "axlle@mail.ru",
		IsEmail:            BoolToBoolPtr(true),
		IsPhone:            BoolToBoolPtr(true),
		Status:             10,
		Password:           "123456",
		RememberToken:      &rememberToken,
		AuthKey:            &authKey,
		AuthToken:          &authToken,
		PasswordResetToken: &passwordResetToken,
		CreatedAt:          &createdAt,
		UpdatedAt:          &updatedAt,
		Roles:              []models.Role{*role},
	}

	err := NewRepo().Create(&user)
	if err != nil {
		log.Printf("Failed to create user %d: %v", n, err.Error())
	}

	log.Println("Database seeded User successfully!")
}
