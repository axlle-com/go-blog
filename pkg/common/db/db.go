package db

import (
	"github.com/axlle-com/blog/pkg/common/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"sync"
)

var (
	once     sync.Once
	instance *gorm.DB
)

func Init(url string) {
	var err error
	once.Do(func() {
		instance, err = gorm.Open(postgres.Open(url), &gorm.Config{})
		if err != nil {
			log.Fatalln(err)
		}
	})
}

func GetDB() *gorm.DB {
	if instance == nil {
		Init(config.GetConfig().DBUrl)
	}
	return instance
}
