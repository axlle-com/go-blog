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

func InitDB(url string) {
	instanceMu.Lock()
	defer instanceMu.Unlock()

	var err error
	instance, err = gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
}

func GetDB() *gorm.DB {
	cfg := config.Config()

	if cfg.IsTest() {
		return GetDBTest()
	}

	instanceMu.Lock()
	defer instanceMu.Unlock()

	if instance == nil {
		var err error
		instance, err = gorm.Open(postgres.Open(cfg.DBUrl()), &gorm.Config{})
		if err != nil {
			logger.Fatalf("[DB][GetDB] Error: %v", err)
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
			logger.Fatalf("[DB][GetDBTest] Error: %v", err)
		}
	}
	return instanceTest //.Debug()
}
