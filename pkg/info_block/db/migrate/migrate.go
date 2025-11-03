package migrate

import (
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/info_block/models"
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
		&models.InfoBlock{},
		&models.InfoBlockHasResource{},
	)

	if err != nil {
		return err
	}

	m.db.Exec(db.HashIndex("info_blocks", "uuid"))
	m.db.Exec(db.HashIndex("info_block_has_resources", "resource_uuid"))

	return nil
}
func (m *migrator) Rollback() error {
	err := m.db.Migrator().DropTable(
		&models.InfoBlock{},
		&models.InfoBlockHasResource{},
	)

	if err != nil {
		return err
	}

	return nil
}
