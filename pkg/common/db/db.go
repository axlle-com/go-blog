package db

import (
	"github.com/axlle-com/blog/pkg/post"
	"github.com/axlle-com/blog/pkg/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func Init(url string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	err = db.AutoMigrate(&post.Post{})
	err = db.AutoMigrate(&user.User{})
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
