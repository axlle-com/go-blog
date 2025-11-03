package db

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/axlle-com/blog/app/models/contract"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

var (
	instance *db
	once     sync.Once
)

type db struct {
	pgsql *gorm.DB
}

func (r *db) PostgreSQL() *gorm.DB { return r.pgsql }

func (r *db) Close() error {
	sqlDB, err := r.pgsql.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func buildGormLogger(cfg contract.Config) glogger.Interface {
	level := glogger.Warn
	if cfg.IsTest() {
		level = glogger.Silent
	}

	switch os.Getenv("GORM_LOG") {
	case "silent":
		level = glogger.Silent
	case "error":
		level = glogger.Error
	case "warn":
		level = glogger.Warn
	case "info", "debug":
		level = glogger.Info
	}
	if os.Getenv("DEBUG") == "1" && !cfg.IsTest() {
		level = glogger.Info
	}

	return glogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		glogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  level,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}

func SetupDB(config contract.Config) (contract.DB, error) {
	var initErr error
	once.Do(func() {
		url := config.DBUrl()
		if config.IsTest() {
			url = config.DBUrlTest()
		}

		dbConn, err := gorm.Open(postgres.Open(url), &gorm.Config{
			Logger: buildGormLogger(config),
		})
		if err != nil {
			initErr = fmt.Errorf("failed to open DB: %w", err)
			return
		}

		if os.Getenv("DEBUG") == "1" && !config.IsTest() {
			dbConn = dbConn.Debug()
		}

		instance = &db{pgsql: dbConn}
	})
	return instance, initErr
}
