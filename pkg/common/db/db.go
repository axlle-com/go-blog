package db

import (
	"github.com/axlle-com/blog/pkg/common/config"
	"github.com/axlle-com/blog/pkg/common/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var (
	instance     *gorm.DB
	instanceTest *gorm.DB
)

func Init(url string) {
	var err error
	instance, err = gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
}

func GetDB() *gorm.DB {
	if config.Config().IsTest() {
		return GetDBTest()
	}

	var err error
	if instance == nil {
		instance, err = gorm.Open(postgres.Open(config.Config().DBUrl()), &gorm.Config{})
		if err != nil {
			logger.Fatal(err)
		}
	}
	return instance
}

func GetDBTest() *gorm.DB {
	var err error
	if instanceTest == nil {
		instanceTest, err = gorm.Open(postgres.Open(config.Config().DBUrlTest()), &gorm.Config{})
		if err != nil {
			log.Fatalln(err)
		}
	}
	return instanceTest //.Debug()
}
