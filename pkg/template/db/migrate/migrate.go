package migrate

import (
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/template/models"
	"gorm.io/gorm"
)

type migrator struct {
	db *gorm.DB
}

func NewMigrator(db *gorm.DB) contract.Migrator {
	return &migrator{db: db}
}

func (m *migrator) Migrate() error {
	err := m.db.AutoMigrate(
		&models.Template{},
	)

	if err != nil {
		return err
	}

	m.db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_templates_theme_name ON templates (theme, name);`)

	return nil
}
func (m *migrator) Rollback() error {
	err := m.db.Migrator().DropTable(
		&models.Template{},
	)

	if err != nil {
		return err
	}

	return nil
}
