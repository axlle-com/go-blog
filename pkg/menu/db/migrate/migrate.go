package migrate

import (
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/menu/models"
	"gorm.io/gorm"
)

type migrator struct {
	db *gorm.DB
}

func NewMigrator(db *gorm.DB) contracts.Migrator {
	return &migrator{db: db}
}

func (m *migrator) Migrate() error {
	err := m.db.AutoMigrate(
		&models.Menu{},
		&models.MenuItem{},
	)

	if err != nil {
		return err
	}

	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_menus_uuid ON menus USING hash (uuid);`)

	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_menu_items_publisher_uuid ON menu_items USING hash (publisher_uuid);`)
	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_menu_items_url ON menu_items USING hash (url);`)

	return nil
}
func (m *migrator) Rollback() error {
	err := m.db.Migrator().DropTable(
		&models.Menu{},
		&models.MenuItem{},
	)

	if err != nil {
		return err
	}

	return nil
}
