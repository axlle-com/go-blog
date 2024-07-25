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

func Init(url string) *gorm.DB {
	var err error
	once.Do(func() {
		instance, err = gorm.Open(postgres.Open(url), &gorm.Config{})
		if err != nil {
			log.Fatalln(err)
		}
	})
	return instance.Debug()
}

func GetDB() *gorm.DB {
	if instance == nil {
		instance = Init(config.GetConfig().DBUrl)
	}
	return instance
}
