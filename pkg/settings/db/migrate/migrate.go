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
	model := &models.Setting{}

	if err := m.db.AutoMigrate(model); err != nil {
		return err
	}

	m.db.Exec(db.UniqueIndex(model.GetTable(), "namespace", "key", "scope"))
	m.db.Exec(db.CompositeIndex(model.GetTable(), "namespace", "scope"))
	m.db.Exec(db.GinIndex(model.GetTable(), "value"))
	m.db.Exec(db.CompositeIndex(model.GetTable(), "sort"))

	return nil
}

func (m *migrator) Rollback() error {
	return m.db.Migrator().DropTable(&models.Setting{})
}
