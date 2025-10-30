package migrate

import (
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/file/models"
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
		&models.File{},
	)
	if err != nil {
		return err
	}

	m.db.Exec(db.CreateHashIndex("files", "uuid"))
	m.db.Exec(db.CreateHashIndex("files", "file"))

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
