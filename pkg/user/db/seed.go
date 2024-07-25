package db

import (
	. "github.com/axlle-com/blog/pkg/common/db"
	. "github.com/axlle-com/blog/pkg/common/models"
	rights "github.com/axlle-com/blog/pkg/rights/repository"
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
		avatar := faker.URL()
		password := faker.Password()
		passwordHash := faker.Password()
		rememberToken := faker.UUIDDigit()
		authKey := faker.UUIDDigit()
		authToken := faker.UUIDHyphenated()
		passwordResetToken := faker.UUIDHyphenated()
		createdAt := time.Now()
		updatedAt := time.Now()

		user := User{
			FirstName:          firstName,
			LastName:           lastName,
			Phone:              &phone,
			Email:              faker.Email(),
			IsEmail:            IntToBoolPtr(),
			IsPhone:            IntToBoolPtr(),
			Status:             int8(rand.Intn(10)),
			Avatar:             &avatar,
			Password:           password,
			PasswordHash:       passwordHash,
			RememberToken:      &rememberToken,
			AuthKey:            &authKey,
			AuthToken:          &authToken,
			PasswordResetToken: &passwordResetToken,
			CreatedAt:          &createdAt,
			UpdatedAt:          &updatedAt,
		}
		err := NewRepository().Create(&user)
		if err != nil {
			log.Printf("Failed to create user %d: %v", i, err.Error())
		}
	}

	firstName := faker.FirstName()
	lastName := faker.LastName()
	phone := faker.Phonenumber()
	avatar := faker.URL()
	passwordHash := faker.Password()
	rememberToken := faker.UUIDDigit()
	authKey := faker.UUIDDigit()
	authToken := faker.UUIDHyphenated()
	passwordResetToken := faker.UUIDHyphenated()
	createdAt := time.Now()
	updatedAt := time.Now()

	role, _ := rights.NewRoleRepository().GetByName("admin")
	user := User{
		FirstName:          firstName,
		LastName:           lastName,
		Phone:              &phone,
		Email:              "axlle@mail.ru",
		IsEmail:            IntToBoolPtr(),
		IsPhone:            IntToBoolPtr(),
		Status:             int8(rand.Intn(10)),
		Avatar:             &avatar,
		Password:           "123456",
		PasswordHash:       passwordHash,
		RememberToken:      &rememberToken,
		AuthKey:            &authKey,
		AuthToken:          &authToken,
		PasswordResetToken: &passwordResetToken,
		CreatedAt:          &createdAt,
		UpdatedAt:          &updatedAt,
		Roles:              []Role{*role},
	}

	err := NewRepository().Create(&user)
	if err != nil {
		log.Printf("Failed to create user %d: %v", n, err.Error())
	}

	log.Println("Database seeded User successfully!")
}
