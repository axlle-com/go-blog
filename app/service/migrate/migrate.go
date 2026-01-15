package migrate

import (
	"github.com/axlle-com/blog/app/models/contract"
	"gorm.io/gorm"
)

type migrator struct {
	db        *gorm.DB
	migrators []contract.Migrator
}

func NewMigrator(db *gorm.DB, arg ...contract.Migrator) contract.Migrator {
	m := &migrator{db: db}
	m.migrators = append(m.migrators, arg...)

	return m
}

func (m *migrator) Migrate() error {
	err := m.db.AutoMigrate(
		&Seed{},
	)
	if err != nil {
		return err
	}

	for _, mig := range m.migrators {
		if err := mig.Migrate(); err != nil {
			return err
		}
	}

	return nil
}

func (m *migrator) Rollback() error {
	err := m.db.Migrator().DropTable(
		&Seed{},
	)

	if err != nil {
		return err
	}

	for _, mig := range m.migrators {
		if err := mig.Rollback(); err != nil {
			return err
		}
	}

	return nil
}
