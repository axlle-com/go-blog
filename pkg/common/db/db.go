package db

import (
	"github.com/axlle-com/blog/pkg/common/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func Init(url string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	err = db.AutoMigrate(&models.Post{})
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
