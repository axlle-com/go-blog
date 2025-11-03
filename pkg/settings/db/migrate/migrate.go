package migrate

import (
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/settings/models"
	"gorm.io/gorm"
)

type migrator struct{ db *gorm.DB }

func NewMigrator(db *gorm.DB) contract.Migrator { return &migrator{db: db} }

func (m *migrator) Migrate() error {
	if err := m.db.AutoMigrate(&models.Setting{}); err != nil {
		return err
	}

	m.db.Exec(db.UniqueIndex("settings", "namespace", "key", "scope"))
	m.db.Exec(db.CompositeIndex("settings", "namespace", "scope"))
	m.db.Exec(db.GinIndex("settings", "value"))
	m.db.Exec(db.CompositeIndex("settings", "sort"))

	return nil
}

func (m *migrator) Rollback() error {
	return m.db.Migrator().DropTable(&models.Setting{})
}
