package migrate

import (
	"github.com/axlle-com/blog/app/models/contracts"
	"gorm.io/gorm"
)

type migrator struct {
	db        *gorm.DB
	migrators []contracts.Migrator
}

func NewMigrator(db *gorm.DB, arg ...contracts.Migrator) contracts.Migrator {
	m := &migrator{db: db}
	m.migrators = append(m.migrators, arg...)
	return m
}

func (m *migrator) Migrate() error {
	for _, mig := range m.migrators {
		if err := mig.Migrate(); err != nil {
			return err
		}
	}
	return nil
}

func (m *migrator) Rollback() error {
	for i := len(m.migrators) - 1; i >= 0; i-- {
		if err := m.migrators[i].Rollback(); err != nil {
			return err
		}
	}
	return nil
}
