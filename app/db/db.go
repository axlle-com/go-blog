package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"

	"github.com/axlle-com/blog/app/models/contracts"
)

var (
	instance *db
	once     sync.Once
)

type db struct {
	gorm *gorm.DB
}

func (r *db) GORM() *gorm.DB {
	return r.gorm
}

func (r *db) Close() error {
	sqlDB, err := r.gorm.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func SetupDB(config contracts.Config) (contracts.DB, error) {
	var initErr error
	once.Do(func() {
		url := config.DBUrl()

		if config.IsTest() {
			url = config.DBUrlTest()
		}

		dbConn, err := gorm.Open(postgres.Open(url), &gorm.Config{})
		if err != nil {
			initErr = fmt.Errorf("failed to open DB: %w", err)
			return
		}
		instance = &db{gorm: dbConn}
	})
	return instance, initErr //.Debug()
}

func GetDB() *gorm.DB {
	if instance == nil {
		panic("db not initialized: call InitDB first")
	}
	return instance.GORM()
}
