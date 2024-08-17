package db

import (
	db "github.com/axlle-com/blog/pkg/common/db"
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
		password := faker.Password()
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

	role, _ := rights.NewRoleRepository().GetByName("admin")
	user := User{
		FirstName:          "Admin",
		LastName:           "Admin",
		Phone:              &phone,
		Email:              "axlle@mail.ru",
		IsEmail:            db.BoolToBoolPtr(true),
		IsPhone:            db.BoolToBoolPtr(true),
		Status:             10,
		Password:           "123456",
		RememberToken:      &rememberToken,
		AuthKey:            &authKey,
		AuthToken:          &authToken,
		PasswordResetToken: &passwordResetToken,
		CreatedAt:          &createdAt,
		UpdatedAt:          &updatedAt,
		Roles:              []Role{*role},
	}

	err := NewRepo().Create(&user)
	if err != nil {
		log.Printf("Failed to create user %d: %v", n, err.Error())
	}

	log.Println("Database seeded User successfully!")
}
