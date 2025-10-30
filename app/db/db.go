package db

import (
	"fmt"
	"sync"

	"github.com/axlle-com/blog/app/models/contract"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	instance *db
	once     sync.Once
)

type db struct {
	pgsql *gorm.DB
}

func (r *db) PostgreSQL() *gorm.DB {
	return r.pgsql
}

func (r *db) Close() error {
	sqlDB, err := r.pgsql.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func SetupDB(config contract.Config) (contract.DB, error) {
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
		instance = &db{pgsql: dbConn}
	})
	return instance, initErr //.Debug()
}
