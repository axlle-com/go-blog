package migrate

import (
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/menu/models"
	"gorm.io/gorm"
)

type migrator struct {
	db *gorm.DB
}

func NewMigrator(db *gorm.DB) contract.Migrator {
	return &migrator{db: db}
}

func (m *migrator) Migrate() error {
	m.db.Exec("CREATE EXTENSION IF NOT EXISTS ltree;")

	err := m.db.AutoMigrate(
		&models.Menu{},
		&models.MenuItem{},
	)

	if err != nil {
		return err
	}

	m.db.Exec(db.HashIndex("menus", "uuid"))
	m.db.Exec(db.HashIndex("menu_items", "publisher_uuid"))
	m.db.Exec(db.HashIndex("menu_items", "url"))
	m.db.Exec(db.LtreeGistIndex("menu_items", "path_ltree"))

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
