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
		instance = &db{pgsql: dbConn}
	})
	return instance, initErr //.Debug()
}
