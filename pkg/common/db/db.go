package db

import (
	"github.com/axlle-com/blog/pkg/common/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var (
	instance *gorm.DB
)

func Init(url string) {
	var err error
	instance, err = gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
}

// GetDB TODO переделать на динамическое
func GetDB() *gorm.DB {
	if instance == nil {
		Init(config.GetConfig().DBUrl)
	}
	return instance.Debug()
}
