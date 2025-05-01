package migrate

import (
	"gorm.io/gorm"

	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/file/models"
)

type migrator struct {
	db *gorm.DB
}

func NewMigrator(db *gorm.DB) contracts.Migrator {
	return &migrator{db: db}
}

func (m *migrator) Migrate() error {
	err := m.db.AutoMigrate(
		&models.File{},
	)
	if err != nil {
		return err
	}

	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_files_uuid ON files USING hash (uuid);`)
	m.db.Exec(`CREATE INDEX IF NOT EXISTS idx_files_file ON files USING hash (file);`)

	return nil
}
func (m *migrator) Rollback() error {
	err := m.db.Migrator().DropTable(
		&models.File{},
	)

	if err != nil {
		return err
	}

	return nil
}
