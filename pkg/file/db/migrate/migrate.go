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
	model := &models.File{}

	err := m.db.AutoMigrate(
		model,
	)
	if err != nil {
		return err
	}

	m.db.Exec(db.HashIndex(model.GetTable(), "uuid"))
	m.db.Exec(db.HashIndex(model.GetTable(), "file"))

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
