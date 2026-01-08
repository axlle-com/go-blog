package testutil

import (
	"sync"
	"testing"

	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/di"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	once     sync.Once
	baseDB   contract.DB
	baseCont *di.Container
	testDB   *gorm.DB
)

func InitTestContainer(t *testing.T) (*gorm.DB, *di.Container) {
	t.Helper()

	once.Do(func() {
		cfg := config.Config()
		cfg.SetTestENV()

		newDB, err := db.SetupDB(cfg)
		if err != nil {
			panic("db not initialized: " + err.Error())
		}

		container := di.NewContainer(cfg, newDB)

		if err := container.Migrator.Migrate(); err != nil {
			panic("migrate failed: " + err.Error())
		}

		baseDB = newDB
		baseCont = container
		testDB = newDB.PostgreSQL()
	})

	require.NotNil(t, testDB)

	return testDB, baseCont
}

func WithTx(t *testing.T) (*gorm.DB, func()) {
	t.Helper()

	if baseDB == nil {
		panic("call InitTestContainer first")
	}

	tx := baseDB.PostgreSQL().Begin()
	if tx.Error != nil {
		t.Fatalf("begin tx: %v", tx.Error)
	}

	cleanup := func() {
		_ = tx.Rollback().Error
	}

	return tx, cleanup
}
