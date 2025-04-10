package db

import (
	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"sync"
)

var (
	instance     *gorm.DB
	instanceTest *gorm.DB

	instanceMu     sync.Mutex
	instanceTestMu sync.Mutex
)

func Init(url string) {
	instanceMu.Lock()
	defer instanceMu.Unlock()

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

	instanceMu.Lock()
	defer instanceMu.Unlock()

	if instance == nil {
		var err error
		instance, err = gorm.Open(postgres.Open(config.Config().DBUrl()), &gorm.Config{})
		if err != nil {
			logger.Fatal(err)
		}
	}
	return instance //.Debug()
}

func GetDBTest() *gorm.DB {
	instanceTestMu.Lock()
	defer instanceTestMu.Unlock()

	if instanceTest == nil {
		var err error
		instanceTest, err = gorm.Open(postgres.Open(config.Config().DBUrlTest()), &gorm.Config{})
		if err != nil {
			log.Fatalln(err)
		}
	}
	return instanceTest //.Debug()
}
