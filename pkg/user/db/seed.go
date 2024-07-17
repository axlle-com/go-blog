package db

import (
	"fmt"
	. "github.com/axlle-com/blog/pkg/common/models"
	. "github.com/axlle-com/blog/pkg/user/repository"
	"log"
	"math/rand"
	"time"
)

func SeedUsers(count int) {
	for i := 1; i <= count; i++ {
		user := User{
			FirstName:          fmt.Sprintf("FirstName%d", i),
			LastName:           fmt.Sprintf("LastName%d", i),
			Patronymic:         fmt.Sprintf("Patronymic%d", i),
			Phone:              fmt.Sprintf("1234567890%d", i),
			Email:              fmt.Sprintf("user%d@example.com", i),
			IsEmail:            uint8(rand.Intn(2)),
			IsPhone:            uint8(rand.Intn(2)),
			Status:             int16(rand.Intn(10)),
			Avatar:             fmt.Sprintf("https://example.com/avatar%d.png", i),
			Password:           fmt.Sprintf("passwordhash%d", i),
			RememberToken:      fmt.Sprintf("remembertoken%d", i),
			AuthKey:            fmt.Sprintf("authkey%d", i),
			AuthToken:          fmt.Sprintf("authtoken%d", i),
			PasswordResetToken: fmt.Sprintf("passwordresettoken%d", i),
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}

		err := NewUserRepository().CreateUser(&user)
		if err != nil {
			log.Printf("Failed to create user %d: %v", i, err.Error())
		}
	}
}
