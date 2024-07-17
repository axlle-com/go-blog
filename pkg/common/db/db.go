package db

import (
	"github.com/axlle-com/blog/pkg/common/config"
	"github.com/axlle-com/blog/pkg/common/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"sync"
)

var (
	once     sync.Once
	instance *gorm.DB
)

func Init(url string) *gorm.DB {
	var db *gorm.DB
	var err error
	once.Do(func() {
		db, err = gorm.Open(postgres.Open(url), &gorm.Config{})
		if err != nil {
			log.Fatalln(err)
		}
		err = db.AutoMigrate(&models.Post{})
		err = db.AutoMigrate(&models.User{})
		if err != nil {
			log.Fatalln(err)
		}
	})
	return db
}

func GetDB() *gorm.DB {
	if instance == nil {
		instance = Init(config.GetConfig().DBUrl)
	}
	return instance
}
