package migrate

import (
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

	m.db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS ux_settings_ns_key_scope ON settings(namespace, key, scope);`)
	m.db.Exec(`CREATE INDEX IF NOT EXISTS ix_settings_ns_scope ON settings(namespace, scope);`)
	m.db.Exec(`CREATE INDEX IF NOT EXISTS ix_settings_value_gin ON settings USING gin(value jsonb_path_ops);`)
	m.db.Exec(`CREATE INDEX IF NOT EXISTS ix_settings_sort ON settings(sort);`)
	return nil
}

func (m *migrator) Rollback() error {
	return m.db.Migrator().DropTable(&models.Setting{})
}
